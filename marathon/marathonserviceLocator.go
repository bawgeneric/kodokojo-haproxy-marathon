package marathon

import (
	"encoding/json"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"io/ioutil"
	"log"
	"net/http"
)

type MarathonServiceLocator struct {
	marathonUrl string
}

func (m MarathonServiceLocator) LocateServiceByType(projectName string, entityType string) ([]commons.Service, bool) {
	url := m.marathonUrl + "/v2/apps?embed=apps.tasks&label=project,componentType"
	if len(projectName) > 0 {
		url = url + ",project==" + projectName
	}
	if len(entityType) > 0 {
		url = url + ",componentType==" + entityType
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		dataJson, _ := ioutil.ReadAll(resp.Body)

		return m.extractServiceFromJson(dataJson)
	}
	return make([]commons.Service, 0), false
}

func (m MarathonServiceLocator) LocateAllService() ([]commons.Service, bool) {
	return m.LocateServiceByType("", "")
}

func (m MarathonServiceLocator) LocateServiceByProject(projectName string) ([]commons.Service, bool) {
	return m.LocateServiceByType(projectName, "")
}

func (m MarathonServiceLocator) extractServiceFromJson(dataJson []byte) (res []commons.Service, succes bool) {

	apps := Apps{}
	json.Unmarshal(dataJson, &apps)
	res = make([]commons.Service, 0)
	for _, app := range apps.Apps {
		projectName, _ := utils.GetAppIdMatchKodokojoProjectName(app.Id)
		if len(projectName) <=0 {
			return 
		}
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
		haProxyHTTPEntries = append(haProxyHTTPEntries, commons.HaProxyEntry{EntryName: app.Labels.ComponentType, Backends: httpBackends})
		haProxySSHEntries = append(haProxySSHEntries, commons.HaProxyEntry{EntryName: app.Labels.ComponentType, Backends: sshBackends})
		res = append(res, commons.Service{ProjectName: projectName, Version: app.Version, LastConfigChangeAt: app.VersionInfo.LastConfigChangeAt, LastScalingAt: app.VersionInfo.LastConfigChangeAt, HaProxySSHEntries: haProxySSHEntries, HaProxyHTTPEntries: haProxyHTTPEntries})
	}
	succes = len(res) > 0
	return
}

func (m *MarathonServiceLocator) isAlive(healthChecks []HealthCheck) bool {
	var res bool = true
	for _, health := range healthChecks {
		if res {
			res = health.Alive
		}
	}
	return res
}

func NewMarathonServiceLocator(marathonUrl string) MarathonServiceLocator {
	return MarathonServiceLocator{marathonUrl}
}
