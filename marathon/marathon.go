package marathon

import (
	"bytes"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"io/ioutil"
	"log"
	"net/http"
)

func RegisterMarathon(configuration commons.Configuration) {
	marathonUrl := configuration.MarathonUrl() + "/v2/eventSubscriptions?callbackUrl=" + configuration.MarathonCallbackUrl()

	resp, err := http.Post(marathonUrl, "application/json", bytes.NewBuffer([]byte(``)))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Registration response :", string(body))
}
