package haproxy

import (
	"encoding/json"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
)

type ProjectConfiguration struct {
	ProjectName string `json:"projectName"`
	SshIp       string `json:"sshIp"`
	SshPort     int    `json:"sshPort"`
}

type ApplicationsState struct {
	configuration                 commons.Configuration
	serviceLocator                utils.ServiceLocator
	haProxyConfigurationGenerator haProxyConfigurator
	haProxyCurrentContext         commons.HaProxyContext
}

func NewApplicationsState(configuration commons.Configuration, serviceLocator utils.ServiceLocator, haProxyConfigurationGenerator haProxyConfigurator, haProxyCurrentContext commons.HaProxyContext) ApplicationsState {
	return ApplicationsState{configuration: configuration, serviceLocator: serviceLocator, haProxyCurrentContext: haProxyCurrentContext, haProxyConfigurationGenerator: haProxyConfigurationGenerator}
}

func (a *ApplicationsState) Start(marathonEventChannel chan commons.MarathonEvent) {
	log.Println("Starting polling of channel")
	go func() {
		for marathonEvent := range marathonEventChannel {
			log.Println("Receive a Marathon Event from Channel", marathonEvent)
			a.handleMarathonEventInHaProxyConfiguration(marathonEvent)
		}
	}()
}

func (a *ApplicationsState) handleMarathonEventInHaProxyConfiguration(marathonEvent commons.MarathonEvent) {
	log.Println("Processing event ", &marathonEvent, "to apply the a new state ?")
	appId := marathonEvent.AppId
	projectName, entityType := utils.GetAppIdMatchKodokojoProjectName(appId)

	if projectName != "" {
		if entityType != "" {
			services := a.serviceLocator.LocateServiceByProject(projectName)
			if len(services) > 0 {
				previous, found := a.findProjectInHaProxyContext(a.haProxyCurrentContext, projectName)
				if found {
					if previous.LastChangedAd.Before(marathonEvent.Timestamp) {
						a.UpdateServicesIfConfigurationChanged(services)
					}
				} else {
					log.Println("Not able to found project for", projectName, "in current HaProcy Context")
					a.UpdateServicesIfConfigurationChanged(services)
				}
			} else {
				log.Println("Not able to found service for", projectName, "and entity type", entityType)
				service := commons.Service{ProjectName: projectName}
				a.UpdateIfConfigurationChanged(service)
			}
		} else {
			log.Println("AppId", appId, "isn't managed by Kodo Kojo")
		}
	} else {
		log.Println("Not abel to extract project name from appId ", appId, "for event ", &marathonEvent)
	}
}

func (a *ApplicationsState) UpdateIfConfigurationChanged(service commons.Service) {
	wrapper := make([]commons.Service, 1)
	wrapper[0] = service
	a.UpdateServicesIfConfigurationChanged(wrapper)
}

func (a *ApplicationsState) UpdateServicesIfConfigurationChanged(services []commons.Service) {
	previous := &(a.haProxyCurrentContext)
	newState := a.haProxyCurrentContext
	for _, service := range services {
		if len(service.ProjectName) > 0 {
			log.Println("Status for project", service.ProjectName, "may had changed.")
			previousProjectPtr, found := a.findProjectInHaProxyContext(a.haProxyCurrentContext, service.ProjectName)
			if found {
				previousProject := &previousProjectPtr

				projectUpdated := *previousProject
				projectUpdated.Version = service.Version
				projectUpdated.LastConfigChangeAt = service.LastConfigChangeAt
				projectUpdated.LastScalingAt = service.LastScalingAt
				projectUpdated.HaProxyHTTPEntries = service.HaProxyHTTPEntries
				projectUpdated.HaProxySSHEntries = service.HaProxySSHEntries

				projects := make([]commons.Project, 0)
				//Adding all other projects
				for _, project := range newState.Projects {
					if project.ProjectName != service.ProjectName {
						projects = append(projects, project)
					}
				}
				projects = append(projects, projectUpdated)
				newState.Projects = projects

				previousConfig := a.haProxyConfigurationGenerator.GenerateConfiguration(*previous)
				newConfig := a.haProxyConfigurationGenerator.GenerateConfiguration(newState)

				previousHash := a.hash(previousConfig)
				newHash := a.hash(newConfig)

				if previousHash != newHash {
					log.Println("Project", service.ProjectName, "update existing project configuration, reloading it.")
					a.haProxyCurrentContext = newState

					a.haProxyConfigurationGenerator.ReloadHaProxyWithConfiguration(newConfig, a.configuration, a.haProxyCurrentContext)
				} else {
					log.Println("Abording change configuration in Ha proxy for project ", service.ProjectName, "No changed detected in configuration.")
				}
			} else {
				log.Println(service.ProjectName, "is a new project.")
				projectConfig := a.getSshConfiguration(service.ProjectName)
				project := commons.Project{ProjectName: service.ProjectName,
					SSHIp:              projectConfig.SshIp,
					SSHPort:            projectConfig.SshPort,
					HaProxySSHEntries:  service.HaProxySSHEntries,
					HaProxyHTTPEntries: service.HaProxyHTTPEntries,
					LastScalingAt:      service.LastScalingAt,
					LastConfigChangeAt: service.LastConfigChangeAt,
					Version:            service.Version}

				a.haProxyCurrentContext.AddProject(project)

				config := a.haProxyConfigurationGenerator.GenerateConfiguration(a.haProxyCurrentContext)
				a.haProxyConfigurationGenerator.ReloadHaProxyWithConfiguration(config, a.configuration, a.haProxyCurrentContext)
			}
		}
	}
}

func (a *ApplicationsState) findProjectInHaProxyContext(haProxyContext commons.HaProxyContext, projectName string) (commons.Project, bool) {
	res := commons.Project{}
	var found bool = false
	for i := 0; i < len(haProxyContext.Projects) && !found; i++ {
		current := haProxyContext.Projects[i]
		if projectName == current.ProjectName {
			res = current
			found = true
		}
	}
	return res, found
}

func (a *ApplicationsState) getSshConfiguration(projectName string) ProjectConfiguration {
	url := a.configuration.MarathonUrl() + "/v2/artifacts/config/" + projectName + ".json"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Unable to retrive configuration for project", projectName)
		return ProjectConfiguration{}
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	projectConfig := ProjectConfiguration{}
	json.Unmarshal(responseBody, &projectConfig)
	return projectConfig
}

func (a *ApplicationsState) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
