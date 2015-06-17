package main

import (
	"fmt"
)

func listenAndReact(rtm *RtmStart, inChan chan *InboundEvent, errChan chan *ErrorEvent, outChan chan *OutboundEvent) {
	for {
		select {
		case event := <-inChan:
			fmt.Println("Received msg ", event.User, event.Channel, event.Text)
		case event := <-errChan:
			fmt.Println("Received err ", event.Type, event.Error)
		}
	}
}
