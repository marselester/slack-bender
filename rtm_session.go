package main

import (
    "fmt"
    "log"
    "io"
    "time"
    "encoding/json"
    "code.google.com/p/go.net/websocket"
)

func main() {
    rtm := requestRtmStart()
    ws_conn, err := connectToMessageServer(rtm.Url)
    if err != nil {
        log.Fatal(err)
    }
    in_chan := make(chan *InboundEvent)
    err_chan := make(chan *ErrorEvent)
    out_chan := make(chan *OutboundEvent)
    go eventReceiverWorker(ws_conn, in_chan, err_chan)
    go eventSenderWorker(ws_conn, out_chan)
    go pingWorker(out_chan)
    for {
        select {
        case event := <- in_chan:
            fmt.Println("Received msg ", event.User, event.Channel, event.Text)
        case event := <- err_chan:
            fmt.Println("Received err ", event.Type, event.Error)
        }
    }
}

func eventReceiverWorker(ws_conn *websocket.Conn, event_chan chan *InboundEvent, err_chan chan *ErrorEvent) {
    var json_blob json.RawMessage
    for {
        json_blob = json.RawMessage{}
        // The server blocks here until a message from the client is received.
        err := websocket.JSON.Receive(ws_conn, &json_blob)
        if err == io.EOF {
            log.Fatal("Disconnected")
        }
        if err != nil {
            log.Fatal(err)
        }

        i, err := parseInboundEvent(json_blob)
        if err != nil && err != ErrUnknownEvent {
            log.Fatal(err)
        }
        switch event := i.(type) {
        case *InboundEvent:
            event_chan <- event
        case *ErrorEvent:
            err_chan <- event
        }
    }
}

func eventSenderWorker(ws_conn *websocket.Conn, c chan *OutboundEvent) {
    for {
        event := <- c
        if err := websocket.JSON.Send(ws_conn, *event); err != nil {
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
        case <- time.After(3 * time.Second):
            c <- event
        }
    }
}
