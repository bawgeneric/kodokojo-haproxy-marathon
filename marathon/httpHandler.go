package marathon

import (
	"encoding/json"
	"fmt"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type marathonEventHandler struct {
	projectName          string
	marathonEventChannel chan commons.MarathonEvent
	// is it usefull ? this data seems not to be used
	EventType string
}

// is it really usefull ? always true, because project name always empty
func (h *marathonEventHandler) Accept(marathonEvent commons.MarathonEvent) bool {
	return strings.HasPrefix(marathonEvent.AppId, h.projectName)
}

func (h *marathonEventHandler) Handle(marathonEvent commons.MarathonEvent) {
	log.Println("Push Marathon event", marathonEvent, "to channel")
	h.marathonEventChannel <- marathonEvent
}

func newMarathonEventHandler(eventType string, marathonEventChannel chan commons.MarathonEvent) *marathonEventHandler {
	res := new(marathonEventHandler)
	res.EventType = eventType
	res.marathonEventChannel = marathonEventChannel
	return res
}

type Server struct {
	configuration         commons.Configuration
	marathonEventHandlers []*marathonEventHandler
}

func NewServer(configuration commons.Configuration, marathonEventChannel chan commons.MarathonEvent) Server {

	marathonEventHandlers := make([]*marathonEventHandler, 4)
	marathonEventHandlers[0] = newMarathonEventHandler("status_update_event", marathonEventChannel)
	marathonEventHandlers[1] = newMarathonEventHandler("health_status_changed_event", marathonEventChannel)
	marathonEventHandlers[2] = newMarathonEventHandler("remove_health_check_event", marathonEventChannel)
	marathonEventHandlers[3] = newMarathonEventHandler("failed_health_check_event", marathonEventChannel)

	return Server{configuration, marathonEventHandlers}
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
	var treated bool = false
	if found {
		for i := 0; i < len(s.marathonEventHandlers) && !treated; i++ {
			handler := s.marathonEventHandlers[i]
			if handler.Accept(event) {
				log.Println("Handler", &handler, "handle event", event)
				treated = true
				handler.Handle(event)
			}
		}
	} else {
		log.Println("Unhandled event", event, ", unable to extract projet name from AppId", event.AppId)
	}

}
