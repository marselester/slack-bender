package main

import (
    "encoding/json"
    "errors"
)

type OutboundEvent struct {
    ID int `json:"id"`
    Type string `json:"type"`
    Channel string `json:"channel"`
    Text string `json:"text"`
}

type TypeEvent struct {
    Type string
}

type InboundEvent struct {
    Type string
    Channel string
    User string
    Text string
}

type SlackError struct {
    Code int
    Msg string
}

type ErrorEvent struct {
    Type string
    Error SlackError
}

var (
    ErrUnknownEvent = errors.New("unknown event")
)

// parseInboundEvent parses a JSON inbound Slack event and returns
// an appropriate event structure.
func parseInboundEvent(jsonBlob json.RawMessage) (interface{}, error) {
    var typeEvent TypeEvent
    if err := json.Unmarshal(jsonBlob, &typeEvent); err != nil {
        return nil, err
    }

    var event interface{}
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
