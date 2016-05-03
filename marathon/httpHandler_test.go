package marathon

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_Handle(t *testing.T) {
	// given
	c := make(chan commons.MarathonEvent, 2)
	defer close(c)
	h := newMarathonEventHandler("status_update_event", c)
	e := commons.MarathonEvent{"status_update_event", "/acme/ci", time.Now(), true, "http://marathon/my-callback"}
	// when
	h.Handle(e)
	// then
	assert.Equal(t, e, <-c)
}

func Test_Accept_true(t *testing.T) {
	// given
	c := make(chan commons.MarathonEvent)
	defer close(c)
	h := newMarathonEventHandler("status_update_event", c)
	h.projectName = "/acme"
	e := commons.MarathonEvent{"status_update_event", "/acme/ci", time.Now(), true, "http://marathon/my-callback"}
	// when
	a := h.Accept(e)
	// then
	assert.True(t, a)
}

//todo explain behavior

//func Test_Accept_false(t *testing.T) {
//	// given
//	c := make(chan commons.MarathonEvent)
//	defer close(c)
//	h := newMarathonEventHandler("status_update_event", c)
//	e := commons.MarathonEvent{"status_update_event", "/acme/ci", time.Now(), true, "http://marathon/my-callback"}
//	// when
//	a := h.Accept(e)
//	// then
//	assert.False(t, a)
//}

func Test_Handler(t *testing.T) {
	// given
	e := commons.MarathonEvent{"status_update_event", "/acme/ci", time.Now(), true, "http://marathon/my-callback"}
	c := make(chan commons.MarathonEvent, 2)
	defer close(c)
	s := new(Server)
	s.marathonEventHandlers = []*marathonEventHandler{newMarathonEventHandler("status_update_event", c)}
	ts := httptest.NewServer(http.HandlerFunc(s.Handler))
	b, _ := json.Marshal(e)
	// when
	res, err := http.Post(ts.URL, "Application/Json", bytes.NewBuffer(b))
	// then
	scan := bufio.NewScanner(res.Body)
	assert.Nil(t, err)
	assert.True(t, scan.Scan())
	assert.Equal(t, "\"OK\"", scan.Text())
	assert.Equal(t, e, <-c)
}
