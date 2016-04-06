package commons_test

import (
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_is_ready_HTTP_project(t *testing.T) {

	httpEntries := make([]commons.HaProxyEntry, 1)
	backends := make([]commons.HaProxyBackEnd, 1)
	backends[0] = commons.HaProxyBackEnd{BackEndHost: "locahost", BackEndPort: 32458}
	httpEntries[0] = commons.HaProxyEntry{EntryName: "ci", Backends: backends}
	project := commons.Project{ProjectName: "acme", HaProxyHTTPEntries: httpEntries}

	assert.True(t, project.IsReady(), "Project must be ready")
	assert.True(t, project.IsHTTPReady(), "Project must be HTTP ready")
	assert.False(t, project.IsSSHReady(), "Project isn't SSH ready")

}

func Test_is_ready_SSH_project(t *testing.T) {

	sshEntries := make([]commons.HaProxyEntry, 1)
	backends := make([]commons.HaProxyBackEnd, 1)
	backends[0] = commons.HaProxyBackEnd{BackEndHost: "locahost", BackEndPort: 420022}
	sshEntries[0] = commons.HaProxyEntry{EntryName: "scm", Backends: backends}
	project := commons.Project{ProjectName: "acme", HaProxySSHEntries: sshEntries}

	assert.True(t, project.IsReady(), "Project must be ready")
	assert.True(t, project.IsSSHReady(), "Project must be SSH ready")
	assert.False(t, project.IsHTTPReady(), "Project isn't HTTP ready")

}

func Test_is_NOT_project(t *testing.T) {

	project := commons.Project{ProjectName: "acme"}

	assert.False(t, project.IsReady(), "Project must not be ready")
	assert.False(t, project.IsHTTPReady(), "Project must not be HTTP ready")
	assert.False(t, project.IsSSHReady(), "Project must not be SSH ready")

}
