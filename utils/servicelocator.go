package utils

import (
	"kodokojo-haproxy-marathon/commons"
	"regexp"
)

const projectNameRegexp string = "/(?P<projectName>.*)/(?P<entityType>.*)"

type ServiceLocator interface {
	LocateAllService() (res []commons.Service)

	LocateServiceByProject(projectName string) (res []commons.Service)
}

type KodoKojoProject struct {
	ProjectName string
	EntityName  string
}

func (k KodoKojoProject) HasEntity() bool {
	return k.EntityName != ""
}

func GetAppIdMatchKodokojoProjectName(appId string) (project KodoKojoProject, founded bool) {

	r := regexp.MustCompile(projectNameRegexp)
	namesRegexp := r.SubexpNames()
	submatch := r.FindAllStringSubmatch(appId, -1)
	if submatch != nil {
		result := submatch[0]

		group := map[string]string{}
		for i, value := range result {
			group[namesRegexp[i]] = value
		}
		if len(group) >= 2 {
			project = KodoKojoProject{group["projectName"], group["entityType"]}
			founded = true
		}
	}
	return
}
