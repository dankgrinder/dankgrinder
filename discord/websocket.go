// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const gatewayURL = "wss://gateway.discord.gg/?encoding=json&v=8"

type WSConn struct {
	underlying *websocket.Conn
	sessionID  string
	rtr        *MessageRouter

	// fatalHandler is used for when a fatal error occurs, not when
	// WSConn.Close() is called.
	fatalHandler func(err error)
	client       Client
	seq          int
	closePinger  chan struct{}
}

type WSConnOpts struct {
	MessageRouter *MessageRouter
	FatalHandler  func(err *websocket.CloseError)
}

func (client Client) NewWSConn(rtr *MessageRouter, fatalHandler func(err error)) (*WSConn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error while establishing websocket connection: %v", err)
	}

	c := WSConn{
		underlying:   conn,
		rtr:          rtr,
		fatalHandler: fatalHandler,
		client:       client,
		closePinger:  make(chan struct{}),
	}

	// Receive hello message
	interval, err := c.readHello()
	if err != nil {
		c.underlying.Close()
		return nil, err
	}

	// Authenticate
	err = c.underlying.WriteJSON(&Event{
		Op: OpcodeIdentify,
		Data: Data{
			ClientState: ClientState{
				HighestLastMessageID:     "0",
				ReadStateVersion:         0,
				UserGuildSettingsVersion: -1,
			},
			Identify: Identify{
				Token: c.client.Token,
				Properties: Properties{
					OS:                "Linux",
					Browser:           "Chrome",
					BrowserUserAgent:  "Chrome/86.0.4240.75",
					BrowserVersion:    "86.0.4240.75",
					Referrer:          "https://discord.com/new",
					ReferringDomain:   "discord.com",
					ReleaseChannel:    "stable",
					ClientBuildNumber: 73683,
				},
				Capabilities: 61,
				Presence: Presence{
					Status: "online",
					Since:  0,
					AFK:    false,
				},
				Compress: false,
			},
		}})
	if err != nil {
		c.underlying.Close()
		return nil, fmt.Errorf("error while sending authentication message: %v", err)
	}

	if err = c.awaitEvent(EventNameReady); err != nil {
		c.underlying.Close()
		return nil, fmt.Errorf("error while awaiting ready message: %v", err)
	}

	go c.ping(interval)
	go c.listen()
	return &c, nil
}

// listen handles incoming websocket messages. This function will not return
// and should therefore be run as a goroutine. Panics if called while WSConn
// instance is already listening.
func (c *WSConn) listen() {
	for {
		_, b, err := c.underlying.ReadMessage()

		if err != nil {
			c.closePinger <- struct{}{}
			c.underlying.Close()
			c.fatalHandler(err)
			break
		}

		var body Event
		if err := json.Unmarshal(b, &body); err != nil {
			// All messages which don't decode properly are likely caused by the
			// data object and are ignored for now.
			continue
		}

		switch body.Op {
		case OpcodeDispatch:
			c.seq = body.Sequence
			if body.Data.SessionID != "" {
				c.sessionID = body.Data.SessionID
			}
			if body.EventName == EventNameMessageCreate ||
				body.EventName == EventNameMessageUpdate {
				go c.rtr.process(body.Data.Message, body.EventName)
			}
		case OpcodeInvalidSession:
			c.fatalHandler(fmt.Errorf("session invalidated"))
			c.Close()
		}
	}
}

// ping periodically sends a heartbeat websocket message. This function will
// not return and should therefore be run as a goroutine. Panics if called
// while WSConn instance is already pinging.
func (c *WSConn) ping(interval time.Duration) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-c.closePinger:
				return
			case <-t.C:
			}
			_ = c.underlying.WriteJSON(&Event{
				Op: OpcodeHeartbeat,
			})
		}
	}()
}

// readHello attempts to read a hello message from the websocket. If the next
// message is not a hello message an error will be returned. Otherwise, the
// heartbeat interval will be returned.
func (c *WSConn) readHello() (time.Duration, error) {
	_, b, err := c.underlying.ReadMessage()
	if err != nil {
		return 0, fmt.Errorf("error while reading message from websocket: %v", err)
	}

	var body Event
	if err := json.Unmarshal(b, &body); err != nil {
		return 0, fmt.Errorf("error while unmarshalling incoming websocket message: %v", err)
	}
	if body.Op != OpcodeHello {
		return 0, fmt.Errorf("unexpected opcode for received websocket message: message is not a hello message")
	}

	if body.Data.HeartbeatInterval <= 0 {
		return 0, fmt.Errorf("unexpected value for heartbeat interval")
	}
	return time.Millisecond * time.Duration(body.Data.HeartbeatInterval), nil
}

// awaitEvent will block until the gateway sends a message with the passed event.
// An error is returned if the next message received from the server is not of
// the correct event name.
func (c *WSConn) awaitEvent(e string) error {
	_, b, err := c.underlying.ReadMessage()
	if err != nil {
		return fmt.Errorf("error while reading message from websocket: %v", err)
	}

	var body Event
	if err = json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("error while unmarshalling incoming websocket message: %v", err)
	}
	if body.EventName != e {
		return fmt.Errorf("unexpected event name for received websocket message: %v, expected %v", body.EventName, e)
	}
	return nil
}

func (c *WSConn) Close() error {
	c.fatalHandler = func(err error) {}
	c.rtr.routes = nil
	c.closePinger <- struct{}{}
	err := c.underlying.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseGoingAway, "going away"),
		time.Now().Add(time.Second*10),
	)
	if err != nil {
		c.underlying.Close()
	}
	return nil
}
