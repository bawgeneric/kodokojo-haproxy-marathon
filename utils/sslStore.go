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

func (s *SslStore) GetPemFileFromSslStore(project string, entityType string) []byte {
	url := s.marathonUrl + apiBaseUrl + "/ssl/" + project + "/" + entityType + "/" + project + "-" + entityType + "-server.pem"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("SSL key", string(res))
		return res
	}
	log.Println("SSL file note found on url", url)
	return nil
}
