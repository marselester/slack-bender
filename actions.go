package main

import (
	"fmt"
	"strings"
)

func listenAndReact(rtm *RtmStart, inChan chan *InboundEvent, errChan chan *ErrorEvent, outChan chan *OutboundEvent) {
	fmt.Println("I am", rtm.Self.Name)
	for {
		select {
		case event := <-inChan:
			go react(rtm, event, outChan)
			fmt.Println("Message:", event.User, event.Channel, event.Text)
		case event := <-errChan:
			fmt.Println("Error:", event.Type, event.Error)
		}
	}
}

// Bot reacts on messages where it was mentioned, for example,
// <@U04KYSKPB>: how are you doing?
func react(rtm *RtmStart, msg *InboundEvent, outChan chan *OutboundEvent) {
	botID := rtm.Self.ID
	isBotMentioned := strings.Contains(msg.Text, botID)
	if !isBotMentioned {
		return
	}

	resp := &OutboundEvent{
		Type:    "message",
		Channel: msg.Channel,
		Text:    "Wazzup",
	}
	outChan <- resp
}
