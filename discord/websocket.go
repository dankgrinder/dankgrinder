// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const gatewayURL = "wss://gateway.discord.gg/?encoding=json&v=8"

const (
	StateListening = 1 << iota
	StatePinging
	StateActive
)

type WSConn struct {
	underlying *websocket.Conn
	sessionID  string
	msgRouter  *MessageRouter
	errHandler func(err error)

	// fatalHandler is used for when a fatal error occurs, not when
	// WSConn.Close() is called.
	fatalHandler func(err *websocket.CloseError)
	token        string
	seq          int
	state        uint8
}

type WSConnOpts struct {
	MessageRouter *MessageRouter
	ErrHandler    func(err error)
	FatalHandler  func(err *websocket.CloseError)
}

func NewWSConn(token string, opts WSConnOpts) (*WSConn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error while establishing websocket connection: %v", err)
	}

	c := WSConn{
		underlying:   conn,
		msgRouter:    opts.MessageRouter,
		errHandler:   opts.ErrHandler,
		fatalHandler: opts.FatalHandler,
		token:        token,
	}

	// Receive hello message
	interval, err := c.readHello()
	if err != nil {
		return nil, err
	}
	go c.pinger(interval)

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
				Token: token,
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
		return nil, fmt.Errorf("error while sending authentication message: %v", err)
	}

	go c.listen()
	c.state |= StateActive
	return &c, nil
}

// listen handles incoming websocket messages. This function will not return
// and should therefore be run as a goroutine. Panics if called while WSConn
// instance is already listening.
func (c *WSConn) listen() {
	if c.state&StateListening == StateListening {
		panic("listen called but WSConn is already listening")
	}
	c.state |= StateListening

	for c.state&StateActive == StateActive {
		_, b, err := c.underlying.ReadMessage()

		if err != nil {
			closeErr, ok := err.(*websocket.CloseError)
			if !ok {
				c.errHandler(fmt.Errorf("error while reading incoming websocket message: %v", err))
				continue
			}
			c.Close()
			if closeErr.Code == websocket.CloseGoingAway {
				if err := c.resume(); err != nil {
					c.fatalHandler(closeErr)
				}
				break
			}
			c.fatalHandler(closeErr)
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
				c.msgRouter.process(body.Data.Message, body.EventName)
			}
		case OpcodeInvalidSession:
			c.Close()
			c.fatalHandler(&websocket.CloseError{Text: "session invalidated"})
			break
		}
	}
}

// pinger periodically sends a heartbeat websocket message. This function will
// not return and should therefore be run as a goroutine. Panics if called
// while WSConn instance is already pinging.
func (c *WSConn) pinger(interval time.Duration) {
	if c.state&StatePinging == StatePinging {
		panic("pinger called but WSConn is already pinging")
	}
	t := time.NewTicker(interval)
	c.state |= StatePinging
	go func() {
		defer t.Stop()
		for c.state&StateActive == StateActive {
			err := c.underlying.WriteJSON(&Event{
				Op: OpcodeHeartbeat,
			})
			if err != nil {
				c.errHandler(fmt.Errorf("error while sending ping: %v", err))
			}
			<-t.C
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

func (c *WSConn) resume() error {
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return fmt.Errorf("error while establishing websocket connection: %v", err)
	}

	*c = WSConn{
		underlying:   conn,
		msgRouter:    c.msgRouter,
		errHandler:   c.errHandler,
		fatalHandler: c.fatalHandler,
		token:        c.token,
		seq:          c.seq,
	}

	interval, err := c.readHello()
	if err != nil {
		return err
	}
	go c.pinger(interval)

	// Authenticate with old session.
	err = c.underlying.WriteJSON(&Event{
		Op: OpcodeResume,
		Data: Data{
			Identify: Identify{
				Token: c.token,
			},
			SessionID: c.sessionID,
			Sequence:  c.seq,
		},
	})
	if err != nil {
		return fmt.Errorf("error while sending resume message: %v", err)
	}

	go c.listen()
	c.state |= StateActive
	return nil
}

func (c *WSConn) Close() error {
	if c.state&StateActive == 0 {
		return fmt.Errorf("already closed")
	}
	c.state = 0
	c.underlying.Close()
	return nil
}
