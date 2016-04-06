package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_haveBackend_should_return_false_when_empty_HaProxyEntry_slice(t *testing.T) {
	// given
	httpEntries := make([]HaProxyEntry, 0)
	project := new(Project)
	// when
	haveBackEnd := project.haveBackend(httpEntries)
	// then
	assert.False(t, haveBackEnd)
}

func Test_haveBackend_should_return_false_when_one_empty_backend_slice(t *testing.T) {
	// given
	backend := *new(HaProxyBackEnd)
	entryWithBackEnd := *new(HaProxyEntry)
	entryWithoutBackEnd := *new(HaProxyEntry)
	entryWithBackEnd.Backends = append(entryWithBackEnd.Backends, backend)
	httpEntries := []HaProxyEntry{entryWithoutBackEnd, entryWithBackEnd}
	project := new(Project)
	// when
	haveBackEnd := project.haveBackend(httpEntries)
	// then
	assert.False(t, haveBackEnd)
}

func Test_haveBackend_should_return_true_when_non_empty_HaProxyEntry_and_backend_slice(t *testing.T) {
	// given
	backend := *new(HaProxyBackEnd)
	firstEntry := *new(HaProxyEntry)
	secondEntry := *new(HaProxyEntry)
	firstEntry.Backends = append(firstEntry.Backends, backend)
	secondEntry.Backends = append(firstEntry.Backends, backend)
	httpEntries := []HaProxyEntry{firstEntry, secondEntry}
	project := new(Project)
	// when
	haveBackEnd := project.haveBackend(httpEntries)
	// then
	assert.True(t, haveBackEnd)
}
