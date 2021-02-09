// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

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
	EventNameMessageUpdate = "MESSAGE_UPDATE"
	EventNameReady         = "READY"
	EventNameResumed       = "RESUMED"
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

type Event struct {
	Op        int    `json:"op"`
	Data      Data   `json:"d,omitempty"`
	Sequence  int    `json:"s,omitempty"`
	EventName string `json:"t,omitempty"`
}

type Data struct {
	Message
	Identify
	ClientState       ClientState `json:"client_state,omitempty"`
	HeartbeatInterval int         `json:"heartbeat_interval,omitempty"`
	SessionID         string      `json:"session_id,omitempty"`
	Sequence          int         `json:"seq,omitempty"` // For sending only
}

type Identify struct {
	Token        string     `json:"token"`
	Properties   Properties `json:"properties"`
	Capabilities int        `json:"capabilities,omitempty"`
	Compress     bool       `json:"compress"`
	Presence     Presence   `json:"presence"`
}

type Presence struct {
	Status     string   `json:"status,omitempty"`
	Since      int      `json:"since,omitempty"`
	Activities []string `json:"activities,omitempty"`
	AFK        bool     `json:"afk,omitempty"`
}

type ClientState struct {
	HighestLastMessageID     string `json:"highest_last_message_id,omitempty"`
	ReadStateVersion         int    `json:"read_state_version,omitempty"`
	UserGuildSettingsVersion int    `json:"user_guild_settings_version,omitempty"`
}

type Properties struct {
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
