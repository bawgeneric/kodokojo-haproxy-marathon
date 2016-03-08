package kodokojo

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"strings"
)

func RegisterMarathon(configuration Configuration) {
	marathonUrl := configuration.marathonUrl + "/v2/eventSubscriptions?callbackUrl=" + configuration.marathonCallbackUrl

	resp, err :=http.Post(marathonUrl,"application/json",bytes.NewBuffer([]byte(``)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Registration response :", string(body))
}

type MarathonEvent struct {
	EventType string `json:"eventType"`
	AppId string `json:"appId"`
	Timestamp string `json:"timestamp"`
	Alive bool `json:"alive"`
	CallbackUrl string `json:"callbackUrl"`
}

func (m *abstractEventHandler) Accept(marathonEvent MarathonEvent) bool {
	return strings.HasPrefix(marathonEvent.AppId, m.projectName)
}

type MarathonEventHandler interface {
	Accept(marathonEvent MarathonEvent) bool
	Handle(marathonEvent MarathonEvent)
}

type abstractEventHandler struct {
	marathonUrl string
	projectName string
}

type SubcribeEventHandler struct {
}

func (h SubcribeEventHandler) Accept(marathonEvent MarathonEvent) bool  {
	return marathonEvent.EventType == "subscribe_event"
}

func (h SubcribeEventHandler) Handle(marathonEvent MarathonEvent)  {
	fmt.Println("New subscription on Url '", marathonEvent.CallbackUrl)
}

type UnSubcribeEventHandler struct {
}

func (h UnSubcribeEventHandler) Accept(marathonEvent MarathonEvent) bool  {
	return marathonEvent.EventType == "unsubscribe_event"
}

func (h UnSubcribeEventHandler) Handle(marathonEvent MarathonEvent)  {
	fmt.Println("New unsubscription on Url '", marathonEvent.CallbackUrl)
}


