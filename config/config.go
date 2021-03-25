// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ShiftStateActive  = "active"
	ShiftStateDormant = "dormant"
)

type Config struct {
	InstancesOpts      []InstanceOpts     `yaml:"instances"`
	Shifts             []Shift            `yaml:"shifts"`
	Features           Features           `yaml:"features"`
	Compat             Compat             `yaml:"compatibility"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
}

type InstanceOpts struct {
	Token              string             `yaml:"token"`
	ChannelID          string             `yaml:"channel_id"`
	IsMaster           bool               `yaml:"is_master"`
	Features           Features           `yaml:"-"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"-"`
	Shifts             []Shift            `yaml:"-"`
}

type override struct {
	InstanceOpts []struct {
		Features           yaml.Node `yaml:"features"`
		SuspicionAvoidance yaml.Node `yaml:"suspicion_avoidance"`
		Shifts             yaml.Node `yaml:"shifts"`
	} `yaml:"instances"`
}

type Compat struct {
	PostmemeOpts         []string `yaml:"postmeme"`
	AllowedSearches      []string `yaml:"allowed_searches"`
	SearchCancel         []string `yaml:"search_cancel"`
	Cooldown             Cooldown `yaml:"cooldown"`
	AwaitResponseTimeout int      `yaml:"await_response_timeout"`
}

type Cooldown struct {
	Beg       int `yaml:"beg"`
	Fish      int `yaml:"fish"`
	Hunt      int `yaml:"hunt"`
	Postmeme  int `yaml:"postmeme"`
	Search    int `yaml:"search"`
	Highlow   int `yaml:"highlow"`
	Blackjack int `yaml:"blackjack"`
	Sell      int `yaml:"sell"`
	Gift      int `yaml:"gift"`
}

type Features struct {
	Commands       Commands        `yaml:"commands"`
	CustomCommands []CustomCommand `yaml:"custom_commands"`
	AutoBuy        AutoBuy         `yaml:"auto_buy"`
	AutoSell       AutoSell        `yaml:"auto_sell"`
	AutoGift       AutoGift        `yaml:"auto_gift"`
	AutoBlackjack  AutoBlackjack   `yaml:"auto_blackjack"`
	AutoShare      AutoShare       `yaml:"auto_share"`
	AutoTidepod    AutoTidepod     `yaml:"auto_tidepod"`
	BalanceCheck   BalanceCheck    `yaml:"balance_check"`
	LogToFile      bool            `yaml:"log_to_file"`
	Debug          bool            `yaml:"debug"`
}

type BalanceCheck struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}

type AutoTidepod struct {
	Enable              bool `yaml:"enable"`
	Interval            int  `yaml:"interval"`
	BuyLifesaverOnDeath bool `yaml:"buy_lifesaver_on_death"`
}

type AutoBlackjack struct {
	Enable            bool                         `yaml:"enable"`
	Priority          bool                         `yaml:"priority"`
	Amount            int                          `yaml:"amount"`
	PauseBelowBalance int                          `yaml:"pause_below_balance"`
	LogicTable        map[string]map[string]string `yaml:"logic_table"`
}

type AutoShare struct {
	Enable         bool `yaml:"enable"`
	MaximumBalance int  `yaml:"maximum_balance"`
	MinimumBalance int  `yaml:"minimum_balance"`
}

type AutoGift struct {
	Enable   bool     `yaml:"enable"`
	Interval int      `yaml:"interval"`
	Items    []string `yaml:"items"`
}

type CustomCommand struct {
	Value             string `yaml:"value"`
	Interval          int    `yaml:"interval"`
	Amount            int    `yaml:"amount"`
	PauseBelowBalance int    `yaml:"pause_below_balance"`
}

type AutoBuy struct {
	FishingPole  bool `yaml:"fishing_pole"`
	HuntingRifle bool `yaml:"hunting_rifle"`
	Laptop       bool `yaml:"laptop"`
}

type AutoSell struct {
	Enable   bool     `yaml:"enable"`
	Interval int      `yaml:"interval"`
	Items    []string `yaml:"items"`
}

type Commands struct {
	Beg      bool `yaml:"beg"`
	Postmeme bool `yaml:"postmeme"`
	Search   bool `yaml:"search"`
	Highlow  bool `yaml:"highlow"`
	Fish     bool `yaml:"fish"`
	Hunt     bool `yaml:"hunt"`
}

type SuspicionAvoidance struct {
	Typing       Typing       `yaml:"typing"`
	MessageDelay MessageDelay `yaml:"message_delay"`
}

type Typing struct {
	Base      int `yaml:"base"`      // A base duration in milliseconds.
	Speed     int `yaml:"speed"`     // Speed in keystrokes per minute.
	Variation int `yaml:"variation"` // A random value in milliseconds from [0,n) added to the base.
}

type MessageDelay struct {
	Base      int `yaml:"base"`      // A base duration in milliseconds.
	Variation int `yaml:"variation"` // A random value in milliseconds from [0,n) added to the base.
}

// Shift indicates an application state (active or dormant) for a duration.
type Shift struct {
	State    string   `yaml:"state"`
	Duration Duration `yaml:"duration"`
}

// Duration is not related to a time.Duration. It is a structure used in a Shift
// type.
type Duration struct {
	Base      int `yaml:"base"`      // A base duration in seconds.
	Variation int `yaml:"variation"` // A random value in seconds from [0,n) added to the base.
}

// Load loads the config from the expected path.
func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file: %v", err)
	}
	defer f.Close()

	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while decoding config: %v", err)
	}

	if _, err = f.Seek(0, 0); err != nil {
		return Config{}, fmt.Errorf("error while seeking back to beginning of file: %v", err)
	}
	var ovr override
	if err = yaml.NewDecoder(f).Decode(&ovr); err != nil {
		return Config{}, fmt.Errorf("error while decoding config override: %v", err)
	}

	if len(cfg.InstancesOpts) != len(ovr.InstanceOpts) {
		panic("amount of instances not equal to the amount of override configs")
	}

	for i, ovrOpts := range ovr.InstanceOpts {
		features := cfg.Features
		sa := cfg.SuspicionAvoidance
		shifts := cfg.Shifts
		if ovrOpts.Features.Kind != 0 {
			if err = ovrOpts.Features.Decode(&features); err != nil {
				return Config{}, fmt.Errorf("instances[%v].features error while decoding: %v", i, err)
			}
		}
		if ovrOpts.SuspicionAvoidance.Kind != 0 {
			if err = ovrOpts.SuspicionAvoidance.Decode(&sa); err != nil {
				return Config{}, fmt.Errorf("instances[%v].suspicion_avoidance error while decoding: %v", i, err)
			}
		}
		if ovrOpts.Shifts.Kind != 0 {
			if err = ovrOpts.Shifts.Decode(&shifts); err != nil {
				return Config{}, fmt.Errorf("instances[%v].shifts error while decoding: %v", i, err)
			}
		}
		cfg.InstancesOpts[i].Features = features
		cfg.InstancesOpts[i].SuspicionAvoidance = sa
		cfg.InstancesOpts[i].Shifts = shifts
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if len(c.InstancesOpts) == 0 {
		return fmt.Errorf("instances: no instances, at least 1 is required")
	}
	if err := validateCompat(c.Compat); err != nil {
		return err
	}
	if err := validateFeatures(c.Features); err != nil {
		return err
	}

	var haveMaster bool
	for i, opts := range c.InstancesOpts {
		if opts.Token == "" {
			return fmt.Errorf("instances[%v]: no token", i)
		}
		if opts.ChannelID == "" {
			return fmt.Errorf("instances[%v]: no channel id", i)
		}
		if !isValidID(opts.ChannelID) {
			return fmt.Errorf("instances[%v]: invalid channel id", i)
		}
		if len(opts.Shifts) == 0 {
			return fmt.Errorf("instances[%v]: no shifts", i)
		}
		if opts.IsMaster {
			if haveMaster {
				return fmt.Errorf("multiple master instances")
			}
			haveMaster = true
		}
		if err := validateFeatures(opts.Features); err != nil {
			return fmt.Errorf("instances[%v]: %v", i, err)
		}
		if err := validateShifts(opts.Shifts); err != nil {
			return fmt.Errorf("instances[%v]: %v", i, err)
		}
	}
	return nil
}

func validateFeatures(features Features) error {
	if features.AutoSell.Enable {
		if features.AutoSell.Interval < 0 {
			return fmt.Errorf("auto-sell interval must be greater than or equal to 0")
		}
		if len(features.AutoSell.Items) == 0 {
			return fmt.Errorf("auto-sell enabled but no items configured")
		}
	}
	if features.AutoGift.Enable {
		if features.AutoGift.Interval < 0 {
			return fmt.Errorf("auto-gift interval must be greater than or equal to 0")
		}
		if len(features.AutoGift.Items) == 0 {
			return fmt.Errorf("auto-gift enabled but no items configured")
		}
	}
	if features.AutoShare.Enable {
		if features.AutoShare.MinimumBalance < 0 {
			return fmt.Errorf("auto-share minimum must be greater than or equal to 0")
		}
		if features.AutoShare.MaximumBalance < 0 {
			return fmt.Errorf("auto-share maximum must be greater than or equal to 0")
		}
		if features.AutoShare.MinimumBalance > features.AutoShare.MaximumBalance {
			return fmt.Errorf("auto-share minumum must be smaller than or equal to maximum")
		}
	}
	if features.AutoTidepod.Enable && features.AutoTidepod.Interval < 0 {
		return fmt.Errorf("auto-tidepod interval must be greater than or equal to 0")
	}
	if features.BalanceCheck.Enable && features.BalanceCheck.Interval <= 0 {
		return fmt.Errorf("balance check interval must be greater than 0")
	}
	if features.AutoBlackjack.Enable {
		if !features.BalanceCheck.Enable {
			return fmt.Errorf("auto-blackjack enabled but balance check disabled")
		}
		if features.AutoBlackjack.Amount < 0 {
			return fmt.Errorf("auto-blackjack amount must be greater than or equal to 0")
		}
		for columnKey, row := range features.AutoBlackjack.LogicTable {
			if columnKey != "A" {
				n, err := strconv.Atoi(columnKey)
				if err != nil || n < 2 || n > 10 {
					return fmt.Errorf("invalid auto-blackjack logic table key: %v", columnKey)
				}
			}
			for rowKey := range row {
				rowKey = strings.Replace(rowKey, "soft", "", -1)
				n, err := strconv.Atoi(rowKey)
				if err != nil || n < 4 || n > 20 {
					return fmt.Errorf("invalid auto-blackjack logic table key: %v", rowKey)
				}
			}
		}
	}

	for i, cmd := range features.CustomCommands {
		if cmd.Value == "" {
			return fmt.Errorf("features.custom_commands[%v].value: no value", i)
		}
		if strings.Contains(cmd.Value, "pls shop") {
			return fmt.Errorf("invalid custom command value: %v, this custom command is disallowed, use auto-gift instead", cmd.Value)
		}
		if strings.Contains(cmd.Value, "pls sell") {
			return fmt.Errorf("invalid custom command value: %v, this custom command is disallowed, use auto-sell instead", cmd.Value)
		}
		if cmd.Amount < 0 {
			return fmt.Errorf("features.custom_commands[%v].amount: value must be greater than or equal to 0", i)
		}
	}
	return nil
}

func validateShifts(shifts []Shift) error {
	for _, shift := range shifts {
		if shift.State != ShiftStateActive && shift.State != ShiftStateDormant {
			return fmt.Errorf("invalid shift state: %v", shift.State)
		}
	}
	return nil
}

func validateCompat(compat Compat) error {
	if len(compat.PostmemeOpts) == 0 {
		return fmt.Errorf("no postmeme compatibility options")
	}
	if len(compat.AllowedSearches) == 0 {
		return fmt.Errorf("no allowed searches")
	}
	if len(compat.SearchCancel) == 0 {
		return fmt.Errorf("no search cancel compatibility options")
	}
	if compat.Cooldown.Postmeme <= 0 {
		return fmt.Errorf("postmeme cooldown must be greater than 0")
	}
	if compat.Cooldown.Hunt <= 0 {
		return fmt.Errorf("hunt cooldown must be greater than 0")
	}
	if compat.Cooldown.Highlow <= 0 {
		return fmt.Errorf("highlow cooldown must be greater than 0")
	}
	if compat.Cooldown.Fish <= 0 {
		return fmt.Errorf("fish cooldown must be greater than 0")
	}
	if compat.Cooldown.Search <= 0 {
		return fmt.Errorf("search cooldown must be greater than 0")
	}
	if compat.Cooldown.Beg <= 0 {
		return fmt.Errorf("beg cooldown must be greater than 0")
	}
	if compat.Cooldown.Gift <= 0 {
		return fmt.Errorf("gift cooldown must be greater than 0")
	}
	if compat.Cooldown.Blackjack <= 0 {
		return fmt.Errorf("blackjack cooldown must be greater than 0")
	}
	if compat.Cooldown.Sell <= 0 {
		return fmt.Errorf("sell cooldown must be greater than 0")
	}
	if compat.AwaitResponseTimeout < 0 {
		return fmt.Errorf("await response timeout must be greater than 0")
	}
	return nil
}

func isValidID(id string) bool {
	return regexp.MustCompile(`^[0-9]+$`).Match([]byte(id))
}
