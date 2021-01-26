// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

const (
	UserFlagDiscordEmployee = 1 << iota
	UserFlagPartneredServerOwner
	UserFlagHypeSquadEvents
	UserFlagBugHunterLevel1
	UserFlagHouseBravery
	UserFlagHouseBrilliance
	UserFlagHouseBalance
	UserFlagEarlySupporter
	UserFlagTeamUser
	UserFlagSystem
	UserFlagBugHunterLevel2
	UserFlagVerifiedBot
	UserFlagEarlyVerifiedBotDeveloper
)

const (
	PremiumTypeNone = iota
	PremiumTypeNitroClassic
	PremiumTypeNitro
)

type User struct {
	// The user's ID.
	ID       string `json:"id"`
	Username string `json:"username"`

	// The user's 4-digit Discord-tag.
	Discriminator string `json:"discriminator"`

	// Whether the user belongs to an OAuth2 application.
	Bot bool `json:"bot,omitempty"`

	// The user's avatar hash.
	Avatar string `json:"avatar"`

	// Whether the user is an official Discord system user (part of the urgent
	// message system).
	System bool `json:"system,omitempty"`

	// Whether the user has two factor authentication enabled on their account.
	MFA bool `json:"mfa_enabled,omitempty"`

	// The user's chosen language option.
	Locale        string `json:"locale,omitempty"`
	VerifiedEmail bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int    `json:"flags,omitempty"`

	// The type of Nitro subscription on a user's account.
	PremiumType int `json:"premium_type,omitempty"`
	PublicFlags int `json:"public_flags,omitempty"`
}
