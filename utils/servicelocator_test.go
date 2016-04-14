package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extract_valid_project_name_and_component_type(t *testing.T) {
	// given
	input := "/acme/ci"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "acme", project.ProjectName)
	assert.Equal(t, "ci", project.EntityName)
	assert.True(t, found)
}

func Test_extract_empty_project_and_empty_component_when_empty_input(t *testing.T) {
	// given
	input := ""
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "", project.ProjectName)
	assert.Equal(t, "", project.EntityName)
	assert.False(t, found)
}

func Test_extract_empty_project_and_empty_component_when_all_delimiters_bad(t *testing.T) {
	// given
	input := "#acme#ci#"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "", project.ProjectName)
	assert.Equal(t, "", project.EntityName)
	assert.False(t, found)
}

func Test_extract_empty_project_and_empty_component_when_first_delimiter_good_but_other_bad(t *testing.T) {
	// given
	input := "/acme#ci#"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "", project.ProjectName)
	assert.Equal(t, "", project.EntityName)
	assert.False(t, found)
}

func Test_extract_empty_project_and_empty_component_when_second_delimiter_good_but_other_bad(t *testing.T) {
	// given
	input := "#acme/ci#"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "", project.ProjectName)
	assert.Equal(t, "", project.EntityName)
	assert.False(t, found)
}

func Test_extract_non_empty_project_but_empty_component_when_empty_component_in_input(t *testing.T) {
	// given
	input := "/acme/"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "acme", project.ProjectName)
	assert.Equal(t, "", project.EntityName)
	assert.True(t, found)
}

func Test_HasEntity_should_return_false_when_empty_entity(t *testing.T) {
	// given
	project := KodoKojoProject{"acme", ""}
	// when
	hasEntity := project.HasEntity()
	// then
	assert.False(t, hasEntity)
}

func Test_HasEntity_should_return_false_when_non_empty_entity(t *testing.T) {
	// given
	project := KodoKojoProject{"acme", "ci"}
	// when
	hasEntity := project.HasEntity()
	// then
	assert.True(t, hasEntity)
}
