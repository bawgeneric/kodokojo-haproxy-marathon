package utils_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"kodokojo-haproxy-marathon/utils"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func Test_GetPemFileFromSslStore_on_bad_url(t *testing.T) {
	// given
	badUrl := "http://bad-url.acme"
	store := utils.NewSslStore(badUrl)
	// when
	data, err := store.GetPemFileFromSslStore("acme", "ci")
	// then
	assert.Empty(t, data)
	assert.NotNil(t, err)
}

func Test_GetPemFileFromSslStore_on_404(t *testing.T) {
	// given
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	store := utils.NewSslStore(ts.URL)
	// when
	data, err := store.GetPemFileFromSslStore("acme", "ci")
	// then
	assert.Empty(t, data)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("SSL file note found on url : %s", ts.URL+"/v2/artifacts/ssl/acme/ci/acme-ci-server.pem"), err.Error())
}

func Test_GetPemFileFromSslStore_with_200_and_result(t *testing.T) {
	// given
	body := strconv.AppendInt(make([]byte, 1), 42, 10)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	store := utils.NewSslStore(ts.URL)
	// when
	data, err := store.GetPemFileFromSslStore("acme", "ci")
	// then
	assert.Equal(t, data, body)
	assert.Nil(t, err)
}
