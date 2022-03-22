package internal

import (
	"github.com/google/uuid"
	"net/http"
)

type Consumer struct {
	id string

	responseWriter http.ResponseWriter
	request        *http.Request
	flusher        http.Flusher

	msg      chan []byte
	doneChan chan interface{}
}

func NewConsumer() (*Consumer, error) {
	return &Consumer{
		id:             uuid.New().String(),
		responseWriter: nil,
		request:        nil,
		flusher:        nil,
		msg:            nil,
		doneChan:       nil,
	}, nil
}
