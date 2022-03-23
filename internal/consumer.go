package internal

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"sync"
)

type Subscription struct {
	id        string
	parent    *EventSource
	stream    chan Event
	errStream chan error
	eventType string
}

func (s *Subscription) Close() {
	s.parent.closeSubscription(s.id)
}

type EventSource struct {
	subscriptions    map[string]*Subscription
	subscriptionsMtx sync.Mutex
	stopChan         chan interface{}
	closed           bool
	unfinishedEvent  *StringEvent
}

func NewEventSource(url string, headers http.Header) (*EventSource, error) {
	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    parsedUrl,
		Header: headers,
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)

	eventSource := &EventSource{
		subscriptions: make(map[string]*Subscription),
		stopChan:      make(chan interface{}),
	}

	go func(response *http.Response, es *EventSource) {
		defer response.Body.Close()

	loop:
		for {
			select {
			case <-es.stopChan:
				break loop
			default:
				b, err := reader.ReadBytes('\n')
				if err != nil && err != io.EOF {
					return
				}

				if len(b) == 0 {
					continue
				}

				// process data
				es.processRaw(b)
			}
		}
	}(resp, eventSource)

	return eventSource, nil
}

func (e *EventSource) processRaw(b []byte) {
	if len(b) == 1 && b[0] == '\n' {
		e.subscriptionsMtx.Lock()
		defer e.subscriptionsMtx.Unlock()

		// previous event is complete
		if e.unfinishedEvent == nil {
			return
		}
		evt := StringEvent{
			Id:    e.unfinishedEvent.Id,
			Event: e.unfinishedEvent.Event,
			Data:  e.unfinishedEvent.Data,
		}
		e.unfinishedEvent = nil
		for _, subscription := range e.subscriptions {
			if subscription.eventType == "" || subscription.eventType == evt.Event {
				subscription.stream <- evt
			}
		}
	}

	payload := strings.TrimRight(string(b), "\n")
	split := strings.SplitN(payload, ":", 2)

	// received comment or heartbeat
	if split[0] == "" {
		return
	}

	if e.unfinishedEvent == nil {
		e.unfinishedEvent = &StringEvent{}
	}

	switch split[0] {
	case "id":
		e.unfinishedEvent.Id = strings.Trim(split[1], " ")
	case "event":
		e.unfinishedEvent.Event = strings.Trim(split[1], " ")
	case "data":
		e.unfinishedEvent.Data = strings.Trim(split[1], " ")
	}
}

func (e *EventSource) Close() {
	if e.closed {
		return
	}

	close(e.stopChan)
	for subId := range e.subscriptions {
		e.closeSubscription(subId)
	}
	e.closed = true
}

func (e *EventSource) Subscribe(eventType string) (*Subscription, error) {
	if e.closed {
		return nil, fmt.Errorf("event source closed")
	}

	sub := &Subscription{
		id:        uuid.New().String(),
		parent:    e,
		stream:    make(chan Event),
		errStream: make(chan error, 1),
		eventType: eventType,
	}

	e.subscriptionsMtx.Lock()
	defer e.subscriptionsMtx.Unlock()

	e.subscriptions[sub.id] = sub

	return sub, nil
}

func (e *EventSource) closeSubscription(id string) bool {
	e.subscriptionsMtx.Lock()
	defer e.subscriptionsMtx.Unlock()

	if sub, ok := e.subscriptions[id]; ok {
		close(sub.stream)
		return true
	}
	return false
}

func (e *EventSource) error(err error) {
	e.subscriptionsMtx.Lock()
	defer e.subscriptionsMtx.Unlock()

	for _, sub := range e.subscriptions {
		sub.errStream <- err
	}

	e.Close()
}
