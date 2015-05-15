package main

import (
	"encoding/json"
	"errors"
)

// OutboundEvent structure is used to send messages.
type OutboundEvent struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// baseEvent is an abstract Slack JSON message.
type baseEvent interface{}

// helperTypeEvent is used to learn event type while parsing, so we can use
// appropriate JSON structure (InboundEvent or ErrorEvent).
type helperTypeEvent struct {
	Type string
}

// InboundEvent represents incoming Slack message.
type InboundEvent struct {
	Type    string
	Channel string
	User    string
	Text    string
}

type slackError struct {
	Code int
	Msg  string
}

// ErrorEvent represents a JSON error we receive from Slack server.
type ErrorEvent struct {
	Type  string
	Error slackError
}

// ErrUnknownEvent is returned when we failed to recognise the event type.
var (
	ErrUnknownEvent = errors.New("unknown event")
)

// parseInboundEvent parses a JSON inbound Slack event and returns
// an appropriate event structure.
func parseInboundEvent(jsonBlob json.RawMessage) (baseEvent, error) {
	var typeEvent helperTypeEvent
	if err := json.Unmarshal(jsonBlob, &typeEvent); err != nil {
		return nil, err
	}

	var event baseEvent
	switch typeEvent.Type {
	case "hello", "message":
		event = new(InboundEvent)
	case "error":
		event = new(ErrorEvent)
	default:
		return nil, ErrUnknownEvent
	}
	if err := json.Unmarshal(jsonBlob, event); err != nil {
		return nil, err
	}

	return event, nil
}
