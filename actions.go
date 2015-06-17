package main

import (
	"fmt"
)

func listenAndReact(rtm *RtmStart, inChan chan *InboundEvent, errChan chan *ErrorEvent, outChan chan *OutboundEvent) {
	fmt.Println("I am", rtm.Self.Name)
	for {
		select {
		case event := <-inChan:
			fmt.Println("Message:", event.User, event.Channel, event.Text)
		case event := <-errChan:
			fmt.Println("Error:", event.Type, event.Error)
		}
	}
}
