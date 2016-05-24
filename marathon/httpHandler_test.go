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
	h := newMarathonEventHandler(c)
	e := commons.MarathonEvent{"status_update_event", "/acme/ci", time.Now(), true, "http://marathon/my-callback"}
	// when
	h.Handle(e)
	// then
	assert.Equal(t, e, <-c)
}

func Test_Handler(t *testing.T) {
	// given
	e := commons.MarathonEvent{"status_update_event", "/acme/ci", time.Now(), true, "http://marathon/my-callback"}
	c := make(chan commons.MarathonEvent, 2)
	defer close(c)
	s := new(Server)
	s.marathonEventHandler = newMarathonEventHandler(c)
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
