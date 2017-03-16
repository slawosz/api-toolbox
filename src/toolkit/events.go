package toolkit

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"net/url"
)

// global events container
//var ec *EventsContainer
//
//func init() {
//	ec = NewEventsContainer()
//}

type Event struct {
	URL      string
	Method   string
	Req      []byte
	Resp     []byte
	RespBody []byte
	UUID     string
}

func (e *Event) SetUUID() {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
	}
	if e.UUID == "" {
		e.UUID = u4.String()
	}
}

type EventsContainer struct {
	EventsList []*Event                       // stores events
	Events     map[*url.URL]map[string]*Event `json:"-"` // stores cached events
}

func NewEventsContainer() *EventsContainer {
	return &EventsContainer{make([]*Event, 0), make(map[*url.URL]map[string]*Event)}
}

func (ec *EventsContainer) HasEvent(req *http.Request) bool {
	method, ok := ec.Events[req.URL]
	if !ok {
		return false
	} else {
		_, ok = method[req.Method]
		return ok
	}
}

func (ec *EventsContainer) AddEvent(req *http.Request, e *Event) {
	ec.EventsList = append(ec.EventsList, e)

	methods, ok := ec.Events[req.URL]
	if !ok {
		methods := make(map[string]*Event)
		ec.Events[req.URL] = methods
	} else {
		methods[req.Method] = e
	}
}
