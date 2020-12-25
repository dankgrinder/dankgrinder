package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const wsURL = "wss://gateway.discord.gg/?encoding=json&v=8"

const (
	OpcodeDispatch = iota
	OpcodeHeartbeat
	OpcodeIdentify
	OpcodePresenceUpdate
	OpcodeVoiceStateUpdate
	OpcodeResume = iota + 1
	OpcodeReconnect
	OpcodeRequestGuildMembers
	OpcodeInvalidSession
	OpcodeHello
	OpcodeHeartbeatACK
)

const (
	EventNameMessageCreate = "MESSAGE_CREATE"
)

const (
	EmbedTypeRich     = "rich"
	EmbedTypeImage    = "image"
	EmbedTypeVideo    = "video"
	EmbedTypeGifVideo = "gifv"
	EmbedTypeArticle  = "article"
	EmbedTypeLink     = "link"
)

const (
	StateListening = 1 << iota
	StatePinging
	StateActive
)

const (
	IntentGuilds = 1 << iota
	IntentGuildMembers
	IntentGuildBans
	IntentGuildEmojis
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping
)

type (
	WSConn struct {
		underlying  *websocket.Conn
		sessionID   string
		chatHandler func(msg Message)
		errHandler  func(err error)

		// fatalHandler is used for when a fatal error occurs, not when
		// WSConn.Close() is called.
		fatalHandler func(err *websocket.CloseError)
		token        string
		seq          int
		state        uint8
	}
	WSConnOpts struct {
		ChatHandler  func(msg Message)
		ErrHandler   func(err error)
		FatalHandler func(err *websocket.CloseError)
	}

	WSMessage struct {
		Op        int    `json:"op"`
		Data      Data   `json:"d,omitempty"`
		Sequence  int    `json:"s,omitempty"`
		EventName string `json:"t,omitempty"`
	}
	Data struct { // TODO: Make identify struct separate (https://discord.com/developers/docs/topics/gateway#identifying)
		Message
		Token             string      `json:"token,omitempty"`
		Properties        Properties  `json:"properties,omitempty"`
		Capabilities      int         `json:"capabilities,omitempty"`
		Presence          Presence    `json:"presence,omitempty"`
		ClientState       ClientState `json:"client_state,omitempty"`
		Compress          bool        `json:"compress,omitempty"`
		HeartbeatInterval int         `json:"heartbeat_interval,omitempty"`
		SessionID         string      `json:"session_id,omitempty"`
		Sequence          int         `json:"seq,omitempty"` // For sending only
	}
	Properties struct {
		OS                     string `json:"os,omitempty"`
		Browser                string `json:"browser,omitempty"`
		Device                 string `json:"device,omitempty"`
		BrowserUserAgent       string `json:"browser_user_agent,omitempty"`
		BrowserVersion         string `json:"browser_version,omitempty"`
		OSVersion              string `json:"os_version,omitempty"`
		Referrer               string `json:"referrer,omitempty"`
		ReferringDomain        string `json:"referring_domain,omitempty"`
		ReferrerCurrent        string `json:"referrer_current,omitempty"`
		ReferringDomainCurrent string `json:"referring_domain_current,omitempty"`
		ReleaseChannel         string `json:"release_channel,omitempty"`
		ClientBuildNumber      int    `json:"client_build_number,omitempty"`
	}
	Presence struct {
		Status     string   `json:"status,omitempty"`
		Since      int      `json:"since,omitempty"`
		Activities []string `json:"activities,omitempty"`
		AFK        bool     `json:"afk,omitempty"`
	}
	ClientState struct {
		HighestLastMessageID     string `json:"highest_last_message_id,omitempty"`
		ReadStateVersion         int    `json:"read_state_version,omitempty"`
		UserGuildSettingsVersion int    `json:"user_guild_settings_version,omitempty"`
	}

	Message struct {
		ID              string    `json:"id,omitempty"`
		ChannelID       string    `json:"channel_id,omitempty"`
		GuildID         string    `json:"guild_id,omitempty"`
		Author          User      `json:"author,omitempty"`
		Content         string    `json:"content,omitempty"`
		Time            time.Time `json:"timestamp,omitempty"`
		EditedTime      time.Time `json:"edited_timestamp,omitempty"`
		TTS             bool      `json:"tts,omitempty"`
		MentionEveryone bool      `json:"mention_everyone,omitempty"`
		Type            int       `json:"type,omitempty"`
		Pinned          bool      `json:"pinned,omitempty"`
		Embeds          []Embed   `json:"embeds,omitempty"`
	}
	User struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
	}
	Embed struct {
		Title       string        `json:"title,omitempty"`
		Type        string        `json:"type,omitempty"`
		Description string        `json:"description,omitempty"`
		URL         string        `json:"url,omitempty"`
		Time        time.Time     `json:"timestamp,omitempty"`
		Color       int           `json:"color,omitempty"`
		Footer      EmbedFooter   `json:"footer,omitempty"`
		Provider    EmbedProvider `json:"provider,omitempty"`
		Author      EmbedAuthor   `json:"author,omitempty"`
		Fields      []EmbedField  `json:"fields,omitempty"`
	}
	EmbedField struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Inline bool   `json:"inline,omitempty"`
	}
	EmbedFooter struct {
		Text         string `json:"text"`
		IconURL      string `json:"icon_url,omitempty"`
		ProxyIconURL string `json:"proxy_icon_url,omitempty"`
	}
	EmbedAuthor struct {
		Name         string `json:"name,omitempty"`
		URL          string `json:"url,omitempty"`
		IconURL      string `json:"icon_url,omitempty"`
		ProxyIconURL string `json:"proxy_icon_url,omitempty"`
	}
	EmbedProvider struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	}
)

func NewWSConn(token string, opts WSConnOpts) (*WSConn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error while establishing websocket connection: %v", err)
	}

	c := WSConn{
		underlying:   conn,
		chatHandler:  opts.ChatHandler,
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
	err = c.underlying.WriteJSON(&WSMessage{
		Op: OpcodeIdentify,
		Data: Data{
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
			Capabilities: 61, // No idea what this means, just mocking the real client.
			Presence: Presence{
				Status: "online",
				Since:  0,
				AFK:    false,
			},
			ClientState: ClientState{
				HighestLastMessageID:     "0",
				ReadStateVersion:         0,
				UserGuildSettingsVersion: -1,
			},
			Compress: false,
		},
	})
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

		var body WSMessage
		if err := json.Unmarshal(b, &body); err != nil {
			c.errHandler(fmt.Errorf("error while unmarshalling incoming websocket message: %v", err))
			continue
		}

		switch body.Op {
		case OpcodeDispatch:
			c.seq = body.Sequence
			if body.Data.SessionID != "" {
				c.sessionID = body.Data.SessionID
			}
			if body.EventName == EventNameMessageCreate {
				c.chatHandler(body.Data.Message)
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
			err := c.underlying.WriteJSON(&WSMessage{
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

	var body WSMessage
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
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("error while establishing websocket connection: %v", err)
	}

	*c = WSConn{
		underlying:   conn,
		chatHandler:  c.chatHandler,
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
	err = c.underlying.WriteJSON(&WSMessage{
		Op: OpcodeResume,
		Data: Data{
			Token:     c.token,
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
