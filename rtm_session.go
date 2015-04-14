package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"code.google.com/p/go.net/websocket"
)

func main() {
	rtm := requestRtmStart()
	wsConn, err := connectToMessageServer(rtm.URL)
	if err != nil {
		log.Fatal(err)
	}
	inChan := make(chan *InboundEvent)
	errChan := make(chan *ErrorEvent)
	outChan := make(chan *OutboundEvent)
	go eventReceiverWorker(wsConn, inChan, errChan)
	go eventSenderWorker(wsConn, outChan)
	go pingWorker(outChan)
	for {
		select {
		case event := <-inChan:
			fmt.Println("Received msg ", event.User, event.Channel, event.Text)
		case event := <-errChan:
			fmt.Println("Received err ", event.Type, event.Error)
		}
	}
}

func eventReceiverWorker(wsConn *websocket.Conn, eventChan chan *InboundEvent, errChan chan *ErrorEvent) {
	var jsonBlob json.RawMessage
	for {
		jsonBlob = json.RawMessage{}
		// The server blocks here until a message from the client is received.
		err := websocket.JSON.Receive(wsConn, &jsonBlob)
		if err == io.EOF {
			log.Fatal("Disconnected")
		}
		if err != nil {
			log.Fatal(err)
		}

		i, err := parseInboundEvent(jsonBlob)
		if err != nil && err != ErrUnknownEvent {
			log.Fatal(err)
		}
		switch event := i.(type) {
		case *InboundEvent:
			eventChan <- event
		case *ErrorEvent:
			errChan <- event
		}
	}
}

func eventSenderWorker(wsConn *websocket.Conn, c chan *OutboundEvent) {
	for {
		event := <-c
		if err := websocket.JSON.Send(wsConn, *event); err != nil {
			log.Fatal(err)
		}
	}
}

func pingWorker(c chan *OutboundEvent) {
	// When there is no other activity clients should send a ping
	// every few seconds.
	event := &OutboundEvent{Type: "ping"}
	for {
		select {
		case <-time.After(3 * time.Second):
			c <- event
		}
	}
}
