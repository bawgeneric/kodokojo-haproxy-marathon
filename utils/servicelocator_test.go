package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_extract_valide_project_name_and_component_type(t *testing.T) {
	input := "/acme/ci"
	projectName, entityType := GetAppIdMatchKodokojoProjectName(input)

	assert.Equal(t, "acme", projectName)
	assert.Equal(t, "ci", entityType)
}
