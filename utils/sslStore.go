package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

const apiBaseUrl string = "/v2/artifacts/"

type SslStore struct {
	marathonUrl string
}

func NewSslStore(marathonUrl string) SslStore {
	return SslStore{marathonUrl: marathonUrl}
}

func (s *SslStore) GetPemFileFromSslStore(project string, entityType string) (res []byte) {
	url := s.marathonUrl + apiBaseUrl + "/ssl/" + project + "/" + entityType + "/" + project + "-" + entityType + "-server.pem"
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	if response.StatusCode == 200 {
		defer response.Body.Close()
		res, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("SSL key", string(res))
	}
	log.Println("SSL file note found on url", url)
	return
}
