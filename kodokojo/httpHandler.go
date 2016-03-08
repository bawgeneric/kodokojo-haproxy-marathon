
package kodokojo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)


type Httphandler struct {
	configuration Configuration
	marathonEventHandlers []MarathonEventHandler
}

func NewHttphandler(configuration Configuration) Httphandler {
	marathonEventHandlers := make([]MarathonEventHandler, 2)
	marathonEventHandlers[0] = SubcribeEventHandler{}
	marathonEventHandlers[1] = UnSubcribeEventHandler{}
	return Httphandler{configuration, marathonEventHandlers}
}

func (h *Httphandler) Start() {
	http.HandleFunc("/callback", h.Handler)
	portStr := fmt.Sprintf(":%d", h.configuration.port)
	http.ListenAndServe(portStr, nil)
}

func (h *Httphandler) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "\"OK\"");
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	event := MarathonEvent{}
	json.Unmarshal(body,&event)
	fmt.Println(event)

	for _,handler := range h.marathonEventHandlers {
		if handler.Accept(event) {
			handler.Handle(event)
		}
	}

}