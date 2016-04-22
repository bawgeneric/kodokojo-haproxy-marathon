package marathon

import (
	"encoding/json"
	"io/ioutil"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"log"
	"net/http"
)

type MarathonServiceLocator struct {
	marathonUrl string
}

func (m MarathonServiceLocator) LocateServiceByProject(projectName string) (res []commons.Service) {
	url := m.marathonUrl + "/v2/apps?embed=apps.tasks&label=project,componentType"
	if len(projectName) > 0 {
		url = url + ",project==" + projectName
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		dataJson, _ := ioutil.ReadAll(resp.Body)
		res = m.ExtractServiceFromJson(dataJson)
	}
	return
}

func (m MarathonServiceLocator) LocateAllService() (res []commons.Service) {
	return m.LocateServiceByProject("")
}

func (m MarathonServiceLocator) ExtractServiceFromJson(dataJson []byte) (res []commons.Service) {

	apps := Apps{}
	json.Unmarshal(dataJson, &apps)
	res = make([]commons.Service, 0)
	for _, app := range apps.Apps {
		project, found := utils.GetAppIdMatchKodokojoProjectName(app.Id)
		if found {
			haProxySSHEntries := make([]commons.HaProxyEntry, 0)
			haProxyHTTPEntries := make([]commons.HaProxyEntry, 0)
			httpBackends := make([]commons.HaProxyBackEnd, 0)
			sshBackends := make([]commons.HaProxyBackEnd, 0)
			for _, task := range app.Tasks {
				if m.isAlive(task.HealthChecks) {
					for i, port := range task.Ports {
						var isHttpPort bool = app.Container.Docker.PortMappings[i].ContainerPort != 22
						backend := commons.HaProxyBackEnd{BackEndHost: task.Host, BackEndPort: port}

						if isHttpPort {
							httpBackends = append(httpBackends, backend)
						} else {
							sshBackends = append(sshBackends, backend)
						}
					}
				}
			}
			if len(httpBackends) > 0 {
				haProxyHTTPEntries = append(haProxyHTTPEntries, commons.HaProxyEntry{EntryName: app.Labels.ComponentType, Backends: httpBackends})
			}
			if len(sshBackends) > 0 {
				haProxySSHEntries = append(haProxySSHEntries, commons.HaProxyEntry{EntryName: app.Labels.ComponentType, Backends: sshBackends})
			}
			res = append(res, commons.Service{ProjectName: project.ProjectName, Version: app.Version, LastConfigChangeAt: app.VersionInfo.LastConfigChangeAt, LastScalingAt: app.VersionInfo.LastConfigChangeAt, HaProxySSHEntries: haProxySSHEntries, HaProxyHTTPEntries: haProxyHTTPEntries})
		}
	}
	return
}

func (m *MarathonServiceLocator) isAlive(healthChecks []HealthCheck) (alive bool) {
	alive = true
	for _, health := range healthChecks {
		if alive {
			alive = health.Alive
		}
	}
	return
}

func NewMarathonServiceLocator(marathonUrl string) (res *MarathonServiceLocator) {
	res = new(MarathonServiceLocator)
	res.marathonUrl = marathonUrl
	return
}
