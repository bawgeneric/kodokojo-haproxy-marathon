package marathon

import (
	"encoding/json"
	"fmt"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"io/ioutil"
	"log"
	"net/http"
)

type marathonEventHandler struct {
	marathonEventChannel chan commons.MarathonEvent
}

func (h *marathonEventHandler) Handle(marathonEvent commons.MarathonEvent) {
	log.Println("Push Marathon event", marathonEvent, "to channel")
	h.marathonEventChannel <- marathonEvent
}

func newMarathonEventHandler(marathonEventChannel chan commons.MarathonEvent) *marathonEventHandler {
	res := new(marathonEventHandler)
	res.marathonEventChannel = marathonEventChannel
	return res
}

type Server struct {
	configuration        commons.Configuration
	marathonEventHandler *marathonEventHandler
}

func NewServer(configuration commons.Configuration, marathonEventChannel chan commons.MarathonEvent) Server {
	return Server{configuration, newMarathonEventHandler(marathonEventChannel)}
}

func (s *Server) Start() {

	// check concurrency
	http.HandleFunc("/callback", s.Handler)
	portStr := fmt.Sprintf(":%d", s.configuration.Port())
	log.Println("Starting to Listen on port", portStr)
	http.ListenAndServe(portStr, nil)
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "\"OK\"")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	event := commons.MarathonEvent{}
	json.Unmarshal(body, &event)
	log.Println(event)

	_, found := utils.GetAppIdMatchKodokojoProjectName(event.AppId)
	if found {
		log.Println("handle event", event)
		s.marathonEventHandler.Handle(event)
	} else {
		log.Println("Unhandled event", event, ", unable to extract projet name from AppId", event.AppId)
	}
}
