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

type Httphandler struct {
	configuration         commons.Configuration
	marathonEventHandlers []MarathonEventHandler
}

func NewHttphandler(configuration commons.Configuration, marathonEventChannel chan commons.MarathonEvent) Httphandler {

	marathonEventHandlers := make([]MarathonEventHandler, 4)
	marathonEventHandlers[0] = newEventTypeEventHandler("status_update_event", marathonEventChannel)
	marathonEventHandlers[1] = newEventTypeEventHandler("health_status_changed_event", marathonEventChannel)
	marathonEventHandlers[2] = newEventTypeEventHandler("remove_health_check_event", marathonEventChannel)
	marathonEventHandlers[3] = newEventTypeEventHandler("failed_health_check_event", marathonEventChannel)

	return Httphandler{configuration, marathonEventHandlers}
}

func (h *Httphandler) Start() {

	http.HandleFunc("/callback", h.Handler)
	portStr := fmt.Sprintf(":%d", h.configuration.Port())
	log.Println("Starting to Listen on port", portStr)
	http.ListenAndServe(portStr, nil)

}

func (h *Httphandler) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "\"OK\"")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	event := commons.MarathonEvent{}
	json.Unmarshal(body, &event)
	log.Println(event)

	projectName, _ := utils.GetAppIdMatchKodokojoProjectName(event.AppId)

	var treated bool = false
	if len(projectName) > 0 {
		for i := 0; i < len(h.marathonEventHandlers) && !treated; i++ {
			handler := h.marathonEventHandlers[i]
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

func (m *abstractEventHandler) Accept(marathonEvent commons.MarathonEvent) bool {
	return strings.HasPrefix(marathonEvent.AppId, m.projectName)
}

type MarathonEventHandler interface {
	Accept(marathonEvent commons.MarathonEvent) bool
	Handle(marathonEvent commons.MarathonEvent)
}

type abstractEventHandler struct {
	projectName          string
	marathonEventChannel chan commons.MarathonEvent
}

func (a *abstractEventHandler) Handle(marathonEvent commons.MarathonEvent) {
	log.Println("Push Marathon event", marathonEvent, "to channel")
	a.marathonEventChannel <- marathonEvent
}

type EventTypeEventHandler struct {
	abstractEventHandler
	EventType string
}

func newEventTypeEventHandler(eventType string, marathonEventChannel chan commons.MarathonEvent) *EventTypeEventHandler {
	res := new(EventTypeEventHandler)
	res.EventType = eventType
	res.marathonEventChannel = marathonEventChannel
	return res
}
