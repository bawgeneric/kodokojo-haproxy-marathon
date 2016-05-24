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
	if found {
		if project.HasEntity() {
			services := a.serviceLocator.LocateAllService()
			if len(services) > 0 {
				a.UpdateServicesIfConfigurationChanged(services)
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
	newState := commons.HaProxyContext{SyslogEntryPoint:a.haProxyCurrentContext.SyslogEntryPoint, ProjectSet:make(map[string]*commons.Project,0)}

	var newConfig string

	httpEntries := make(map[string][]commons.HaProxyEntry,0)
	sshEntries := make(map[string][]commons.HaProxyEntry,0)

	for _, service := range services {
		projectName := service.ProjectName
		if len(projectName) > 0 {
			log.Println("Status for project", projectName, "may had changed.")

			projectConfig := a.getSshConfiguration(projectName)
			project := commons.Project{ProjectName: projectName,
				SSHPort:            projectConfig.SshPort,
				LastScalingAt:      service.LastScalingAt,
				LastConfigChangeAt: service.LastConfigChangeAt,
				Version:            service.Version}

			httpEntries[projectName] = append(httpEntries[projectName], service.HaProxyHTTPEntries...)
			sshEntries[projectName] = append(sshEntries[projectName], service.HaProxySSHEntries...)
			newState.AddProject(project)
		}
	}

	log.Println("Map of HTTPEntries content", &httpEntries)
	log.Println("Map of SSHEntries content", &sshEntries)
	log.Println("NewState =", &newState)
	a.haProxyCurrentContext = newState
	valueAdded := false
	for projectName,project := range a.haProxyCurrentContext.ProjectSet {
		log.Println("Initial HTTP value", &(project))
		log.Println("Adding haEntry for project", projectName, httpEntries[projectName], sshEntries[projectName])
		if !valueAdded {
			valueAdded = len(httpEntries[projectName]) > 0 || len(sshEntries[projectName]) > 0
		}
		project.HaProxyHTTPEntries = httpEntries[projectName]
		project.HaProxySSHEntries = sshEntries[projectName]
	}

	if valueAdded {
		newConfig = a.haProxyConfigurationGenerator.GenerateConfiguration(a.haProxyCurrentContext)
		a.haProxyConfigurationGenerator.ReloadHaProxyWithConfiguration(newConfig, a.configuration, a.haProxyCurrentContext)
	} else {
		log.Println("No backend entries availables, abording reload of Haproxy configuration.")
	}

}

func (a *ApplicationsState) findProjectInHaProxyContext(haProxyContext commons.HaProxyContext, projectName string) (res *commons.Project) {
	res = haProxyContext.ProjectSet[projectName]
	return res
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
