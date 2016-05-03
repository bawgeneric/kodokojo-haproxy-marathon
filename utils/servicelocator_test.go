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
	assert.Empty(t, project.ProjectName)
	assert.Empty(t, project.EntityName)
	assert.False(t, found)
}

func Test_extract_empty_project_and_empty_component_when_all_delimiters_bad(t *testing.T) {
	// given
	input := "#acme#ci#"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Empty(t, project.ProjectName)
	assert.Empty(t, project.EntityName)
	assert.False(t, found)
}

func Test_extract_empty_project_and_empty_component_when_first_delimiter_good_but_other_bad(t *testing.T) {
	// given
	input := "/acme#ci#"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Empty(t, project.ProjectName)
	assert.Empty(t, project.EntityName)
	assert.False(t, found)
}

func Test_extract_empty_project_and_empty_component_when_second_delimiter_good_but_other_bad(t *testing.T) {
	// given
	input := "#acme/ci#"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Empty(t, project.ProjectName)
	assert.Empty(t, project.EntityName)
	assert.False(t, found)
}

func Test_extract_non_empty_project_but_empty_component_when_empty_component_in_input(t *testing.T) {
	// given
	input := "/acme/"
	// when
	project, found := GetAppIdMatchKodokojoProjectName(input)
	// then
	assert.Equal(t, "acme", project.ProjectName)
	assert.Empty(t, project.EntityName)
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

func Test_buildKodokojoProject_on_empty_slice(t *testing.T) {
	// given
	slice := []string{}
	// when
	project, found := buildKodokojoProject(slice)
	// then
	assert.False(t, found)
	assert.Empty(t, project.EntityName)
	assert.Empty(t, project.ProjectName)
}

func Test_buildKodokojoProject_on_slice_with_one_element(t *testing.T) {
	// given
	slice := []string{"/acme/ci"}
	// when
	project, found := buildKodokojoProject(slice)
	// then
	assert.False(t, found)
	assert.Empty(t, project.EntityName)
	assert.Empty(t, project.ProjectName)
}

func Test_buildKodokojoProject_on_slice_with_two_element(t *testing.T) {
	// given
	slice := []string{"/acme/ci", "acme"}
	// when
	project, found := buildKodokojoProject(slice)
	// then
	assert.True(t, found)
	assert.Empty(t, project.EntityName)
	assert.Equal(t, "acme", project.ProjectName)
}

func Test_buildKodokojoProject_on_slice_with_three_element(t *testing.T) {
	// given
	slice := []string{"/acme/ci", "acme", "ci"}
	// when
	project, found := buildKodokojoProject(slice)
	// then
	assert.True(t, found)
	assert.Equal(t, "ci", project.EntityName)
	assert.Equal(t, "acme", project.ProjectName)
}

func Benchmark_getAppIdMatchKodokojoProjectName(b *testing.B) {
	// given
	for i := 0; i < b.N; i++ {
		GetAppIdMatchKodokojoProjectName("/acme/ci")
	}
}
