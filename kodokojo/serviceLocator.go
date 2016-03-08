package kodokojo

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type marathonServiceLocator struct {
	marathonUrl string
}

func (m *marathonServiceLocator) LocateServices(projectName string, entityType string) Services {
	resp, err := http.Get(m.marathonUrl + "/v2/apps?embed=apps.tasks&label=project==" + projectName)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		dataJson, _ := ioutil.ReadAll(resp.Body)

		return m.extractServiceFromJson(projectName, entityType, dataJson)
	}
	return Services{}
}

func (m *marathonServiceLocator) extractServiceFromJson(projectName string, entityType string, dataJson []byte) Services {

	apps := Apps{}
	json.Unmarshal(dataJson, &apps)
	for _, app := range apps.Apps {

	haProxySSHEntries := make([]HaProxyEntry, 0)
	haProxyHTTPEntries := make([]HaProxyEntry, 0)
		httpBackends := make([]HaProxyBackEnd, 0)
		sshBackends := make([]HaProxyBackEnd, 0)
		for _, task := range app.Tasks {
			if m.isAlive(task.HealthChecks) {
				for i, port := range task.Ports {
					var isHttpPort bool = app.Container.Docker.PortMappings[i].ContainerPort != 22
					backend := HaProxyBackEnd{BackEndHost:task.Host, BackEndPort:port}
					if (isHttpPort) {
						httpBackends = append(httpBackends, backend)
					} else {
						sshBackends = append(sshBackends, backend)
					}
				}
			}
		}
		haProxyHTTPEntries = append(haProxyHTTPEntries, HaProxyEntry{EntryName:app.Labels.ComponentType, Backends:httpBackends})
		haProxySSHEntries = append(haProxySSHEntries, HaProxyEntry{EntryName:app.Labels.ComponentType, Backends:sshBackends})
		return Services{ProjectName:projectName, HaProxySSHEntries:haProxySSHEntries, HaProxyHTTPEntries:haProxyHTTPEntries}
	}

	return Services{}
}


func (m *marathonServiceLocator) isAlive(healthChecks []HealthCheck) bool {
	var res bool = true
	for _, health := range healthChecks {
		if res {
			res = health.Alive
		}
	}
	return res
}

func NewMarathonServiceLocator(marathonUrl string) marathonServiceLocator {
	return marathonServiceLocator{marathonUrl}
}

type Services struct {
	ProjectName        string
	HaProxyHTTPEntries []HaProxyEntry
	HaProxySSHEntries  []HaProxyEntry
}

type Apps struct {
	Apps []App `json:"apps"`
}

type App struct {
	Id        string `json:"id"`
	Container Container `json:"container"`
	Labels    Labels `json:"labels"`
	Tasks     []Tasks `json:"tasks"`
}

type Container struct {
	Docker Docker `json:"docker"`
}

type Docker struct {
	PortMappings []PortMapping `json:"portMappings"`
}

type PortMapping struct {
	ContainerPort int `json:"containerPort"`
}

type Labels struct {
	Project       string `json:"project"`
	ComponentType string `json:"componentType"`
	Component     string `json:"component"`
}

type Tasks struct {
	Host         string `json:"host"`
	Ports        []int `json:"ports"`
	HealthChecks []HealthCheck `json:"healthCheckResults"`
}

type HealthCheck struct {
	Alive bool `json:"alive"`
}


