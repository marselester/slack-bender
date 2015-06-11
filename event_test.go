package main

import (
	"encoding/json"
	"testing"
)

func TestInboundEvent(t *testing.T) {
	jsonBlob := json.RawMessage(`{
        "type": "message",
        "ts": "1358878749.000002",
        "user": "U023BECGF",
        "text": "Hello"
    }`)
	expectedEvent := InboundEvent{
		Type: "message",
		Text: "Hello",
		User: "U023BECGF",
	}
	i, _ := parseInboundEvent(jsonBlob)
	event := i.(*InboundEvent)
	if expectedEvent != *event {
		t.Errorf("Expected %#v got %#v", expectedEvent, event)
	}
}

func TestErrorEvent(t *testing.T) {
	jsonBlob := json.RawMessage(`{
        "type": "error",
        "error": {
            "code": 1,
            "msg": "Socket URL has expired"
        }
    }`)
	expectedEvent := ErrorEvent{
		Type: "error",
		Error: slackError{
			Code: 1,
			Msg:  "Socket URL has expired",
		},
	}
	i, _ := parseInboundEvent(jsonBlob)
	event := i.(*ErrorEvent)
	if expectedEvent != *event {
		t.Errorf("Expected %#v got %#v", expectedEvent, event)
	}
}
