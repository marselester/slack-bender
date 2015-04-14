/*
To begin a RTM session make an authenticated call to the `rtm.start`
API method. This provides an initial set of team metadata and
a message server WebSocket URL. Once you have connected
to the message server it will provide a stream of events.
*/
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"code.google.com/p/go.net/websocket"
)

const rtmStartURL = "https://slack.com/api/rtm.start?token=YOUR-TOKEN"

type RtmStart struct {
	URL string
}

var httpClient = &http.Client{}

func requestRtmStart() *RtmStart {
	rtmStart := &RtmStart{}
	resp, err := httpClient.Get(rtmStartURL)
	if err != nil {
		return rtmStart
	}
	defer resp.Body.Close()

	body := make([]byte, 0, 1024)
	chunk := make([]byte, 512)
	for {
		rcvCount, err := resp.Body.Read(chunk)
		for _, val := range chunk[:rcvCount] {
			body = append(body, val)
		}
		if err == io.EOF {
			break
		}
	}

	err = json.Unmarshal(body, rtmStart)
	return rtmStart
}

func connectToMessageServer(wsURL string) (*websocket.Conn, error) {
	urlWithPort, err := addPortToURL(wsURL)
	if err != nil {
		return nil, err
	}
	protocol := ""
	origin := "http://localhost/"
	wsConn, err := websocket.Dial(urlWithPort, protocol, origin)
	return wsConn, err
}

func addPortToURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	sslPort := ":443"
	u.Host = u.Host + sslPort
	return u.String(), nil
}
