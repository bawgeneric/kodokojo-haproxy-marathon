package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const apiBaseUrl string = "/v2/artifacts"

type SslStore struct {
	marathonUrl string
}

func NewSslStore(marathonUrl string) SslStore {
	return SslStore{marathonUrl: marathonUrl}
}

func (s *SslStore) GetPemFileFromSslStore(project string, entityType string) (result []byte, err error) {
	url := s.marathonUrl + apiBaseUrl + "/ssl/" + project + "/" + entityType + "/" + project + "-" + entityType + "-server.pem"
	response, e := http.Get(url)
	if e != nil {
		log.Println(e)
		err = e
	} else if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		result, err = ioutil.ReadAll(response.Body)
		return
	} else {
		err = errors.New(fmt.Sprintf("SSL file note found on url : %s", url))
	}
	return
}
