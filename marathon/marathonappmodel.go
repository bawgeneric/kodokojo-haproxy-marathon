package marathon

import "time"

type Apps struct {
	Apps []App `json:"apps"`
}

type App struct {
	Id          string      `json:"id"`
	Container   Container   `json:"container"`
	Labels      Labels      `json:"labels"`
	Tasks       []Tasks     `json:"tasks"`
	Version     time.Time   `json:"version"`
	VersionInfo VersionInfo `json:"versionInfo"`
}

type VersionInfo struct {
	LastConfigChangeAt time.Time `json:"lastConfigChangeAt"`
	LastScalingAt      time.Time `json:"lastScalingAt"`
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
	Host         string        `json:"host"`
	Ports        []int         `json:"ports"`
	HealthChecks []HealthCheck `json:"healthCheckResults"`
	Version      time.Time     `json:"version"`
}

type HealthCheck struct {
	Alive bool `json:"alive"`
}
