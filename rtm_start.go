/*
To begin a RTM session make an authenticated call to the `rtm.start`
API method. This provides an initial set of team metadata and
a message server WebSocket URL. Once you have connected
to the message server it will provide a stream of events.
*/
package main

import (
    "code.google.com/p/go.net/websocket"
    "net/http"
    "net/url"
    "io"
    "encoding/json"
)

const RTM_START_URL = "https://slack.com/api/rtm.start?token=YOUR-TOKEN"

type RtmStart struct {
    Url string
}

var http_client *http.Client = &http.Client{}

func requestRtmStart() *RtmStart {
    rtm_start := &RtmStart{}
    resp, err := http_client.Get(RTM_START_URL)
    if err != nil {
        return rtm_start
    }
    defer resp.Body.Close()

    body := make([]byte, 0, 1024)
    chunk := make([]byte, 512)
    for {
        rcv_count, err := resp.Body.Read(chunk)
        for _, val := range(chunk[:rcv_count]) {
            body = append(body, val)
        }
        if err == io.EOF {
            break
        }
    }

    err = json.Unmarshal(body, rtm_start)
    return rtm_start
}

func connectToMessageServer(ws_url string) (*websocket.Conn, error) {
    url_, err := addPortToUrl(ws_url)
    if err != nil {
        return nil, err
    }
    protocol := ""
    origin := "http://localhost/"
    ws_conn, err := websocket.Dial(url_, protocol, origin)
    return ws_conn, err
}

func addPortToUrl(url_ string) (string, error) {
    u, err := url.Parse(url_)
    if err != nil {
        return "", err
    }
    ssl_port := ":443"
    u.Host = u.Host + ssl_port
    return u.String(), nil
}
