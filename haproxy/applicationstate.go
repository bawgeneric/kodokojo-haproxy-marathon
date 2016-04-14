package haproxy

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"kodokojo-haproxy-marathon/commons"
	"kodokojo-haproxy-marathon/utils"
	"log"
	"net/http"
)

type ProjectConfiguration struct {
	ProjectName string `json:"projectName"`
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
	project, found := utils.GetAppIdMatchKodokojoProjectName(appId)
	if !found {
		if !project.HasEntity() {
			services := a.serviceLocator.LocateServiceByProject(project.ProjectName)
			if len(services) > 0 {
				previous, found := a.findProjectInHaProxyContext(a.haProxyCurrentContext, project.ProjectName)
				if found {
					if previous.LastChangedAd.Before(marathonEvent.Timestamp) {
						a.UpdateServicesIfConfigurationChanged(services)
					}
				} else {
					log.Println("Not able to found project for", project.ProjectName, "in current HaProcy Context")
					a.UpdateServicesIfConfigurationChanged(services)
				}
			} else {
				log.Println("Not able to found service for", project.ProjectName, "Updating configuration for this project.")
				service := commons.Service{ProjectName: project.ProjectName}
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
	newState := a.haProxyCurrentContext

	var newConfig string

	httpEntries := make(map[string][]commons.HaProxyEntry)
	sshEntries := make(map[string][]commons.HaProxyEntry)

	for _, service := range services {
		if len(service.ProjectName) > 0 {
			log.Println("Status for project", service.ProjectName, "may had changed.")
			previousProject, found := a.findProjectInHaProxyContext(a.haProxyCurrentContext, service.ProjectName)
			if !found {
				previousProjectInNewState, foundInNewState := a.findProjectInHaProxyContext(newState, service.ProjectName)
				previousProject = previousProjectInNewState
				found = foundInNewState
			}
			if found {
				previousProject.Version = service.Version
				previousProject.LastConfigChangeAt = service.LastConfigChangeAt
				previousProject.LastScalingAt = service.LastScalingAt
				if len(service.HaProxyHTTPEntries) > 0 {
					httpEntries[service.ProjectName] = append(service.HaProxyHTTPEntries, httpEntries[service.ProjectName]...)
				}
				if len(service.HaProxySSHEntries) > 0 {
					sshEntries[service.ProjectName] = append(service.HaProxySSHEntries, sshEntries[service.ProjectName]...)
				}

			} else {
				log.Println(service.ProjectName, "is a new project.")
				projectConfig := a.getSshConfiguration(service.ProjectName)
				project := commons.Project{ProjectName: service.ProjectName,
					SSHPort:            projectConfig.SshPort,
					LastScalingAt:      service.LastScalingAt,
					LastConfigChangeAt: service.LastConfigChangeAt,
					Version:            service.Version}

				httpEntries[service.ProjectName] = append(httpEntries[service.ProjectName], service.HaProxyHTTPEntries...)
				sshEntries[service.ProjectName] = append(sshEntries[service.ProjectName], service.HaProxySSHEntries...)
				newState.AddProject(project)
			}
		}
	}

	log.Println("Map of HTTPEntries content", &httpEntries)
	log.Println("Map of SSHEntries content", &sshEntries)
	log.Println("NewState =", &newState)
	a.haProxyCurrentContext = newState
	for i := 0; i < len(a.haProxyCurrentContext.Projects); i++ {
		log.Println("Initial HTTP value", &(a.haProxyCurrentContext.Projects[i]))
		log.Println("Adding haEntry for project", a.haProxyCurrentContext.Projects[i].ProjectName, httpEntries[a.haProxyCurrentContext.Projects[i].ProjectName], sshEntries[a.haProxyCurrentContext.Projects[i].ProjectName])
		a.haProxyCurrentContext.Projects[i].HaProxyHTTPEntries = httpEntries[a.haProxyCurrentContext.Projects[i].ProjectName]
		a.haProxyCurrentContext.Projects[i].HaProxySSHEntries = sshEntries[a.haProxyCurrentContext.Projects[i].ProjectName]
	}
	newConfig = a.haProxyConfigurationGenerator.GenerateConfiguration(a.haProxyCurrentContext)
	a.haProxyConfigurationGenerator.ReloadHaProxyWithConfiguration(newConfig, a.configuration, a.haProxyCurrentContext)

}

func (a *ApplicationsState) findProjectInHaProxyContext(haProxyContext commons.HaProxyContext, projectName string) (res *commons.Project, found bool) {
	found = false
	for i := 0; i < len(haProxyContext.Projects) && !found; i++ {
		if projectName == haProxyContext.Projects[i].ProjectName {
			res = &(haProxyContext.Projects[i])
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
