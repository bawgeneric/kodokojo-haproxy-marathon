package utils

import (
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"regexp"
)

const projectNameRegexp string = "/(.*)/(.*)"

type ServiceLocator interface {
	LocateAllService() (res []commons.Service)

	LocateServiceByProject(projectName string) (res []commons.Service)
}

// Simple structure containing the project name and the entity name
type KodoKojoProject struct {
	ProjectName string
	EntityName  string
}

func (k KodoKojoProject) HasEntity() bool {
	return k.EntityName != ""
}

// extract a KodoKojoProject from the given appId
// return one KodoKojoProject and true if at least the project name is found from the appId, else false
func GetAppIdMatchKodokojoProjectName(appId string) (project KodoKojoProject, found bool) {
	r := regexp.MustCompile(projectNameRegexp)
	submatch := r.FindAllStringSubmatch(appId, -1)
	if submatch != nil {
		project, found = buildKodokojoProject(submatch[0])
	}
	return
}

func buildKodokojoProject(result []string) (project KodoKojoProject, found bool) {
	// the first value of result is the whole expression
	// the second (if it exist) is the project name
	// the third (if it exist again) is the entity name
	if len(result) > 1 {
		found = true
		project.ProjectName = result[1]
	}
	if len(result) > 2 {
		project.EntityName = result[2]
	}
	return
}
