package main

import (
    "testing"
    "encoding/json"
)

func TestInboundEvent(t *testing.T) {
    json_blob := json.RawMessage(`{
        "type": "message",
        "ts": "1358878749.000002",
        "user": "U023BECGF",
        "text": "Hello"
    }`)
    expected_event := InboundEvent{
        Type: "message",
        Text: "Hello",
        User: "U023BECGF",
    }
    i, _ := parseInboundEvent(json_blob)
    event := i.(*InboundEvent)
    if expected_event != *event {
        t.Errorf("Expected %s got %s", expected_event, event)
    }
}

func TestErrorEvent(t *testing.T) {
    json_blob := json.RawMessage(`{
        "type": "error",
        "error": {
            "code": 1,
            "msg": "Socket URL has expired"
        }
    }`)
    expected_event := ErrorEvent{
        Type: "error",
        Error: SlackError{
            Code: 1,
            Msg: "Socket URL has expired",
        },
    }
    i, _ := parseInboundEvent(json_blob)
    event := i.(*ErrorEvent)
    if expected_event != *event {
        t.Errorf("Expected %s got %s", expected_event, event)
    }
}
