// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	ShiftStateActive  = "active"
	ShiftStateDormant = "dormant"
)

type Config struct {
	InstancesOpts      []InstanceOpts     `yaml:"instances"`
	Features           Features           `yaml:"features"`
	Compat             Compat             `yaml:"compatibility"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
}

type InstanceOpts struct {
	Token     string  `yaml:"token"`
	ChannelID string  `yaml:"channel_id"`
	Shifts    []Shift `yaml:"shifts"`
}

type Compat struct {
	PostmemeOpts         []string `yaml:"postmeme"`
	AllowedSearches      []string `yaml:"allowed_searches"`
	SearchCancel         []string `yaml:"search_cancel"`
	Cooldown             Cooldown `yaml:"cooldown"`
	AwaitResponseTimeout int      `yaml:"await_response_timeout"`
}

type Cooldown struct {
	Beg      int `yaml:"beg"`
	Fish     int `yaml:"fish"`
	Hunt     int `yaml:"hunt"`
	Postmeme int `yaml:"postmeme"`
	Search   int `yaml:"search"`
	Highlow  int `yaml:"highlow"`
	Margin   int `yaml:"margin"`
}

type Features struct {
	Commands       Commands        `yaml:"commands"`
	CustomCommands []CustomCommand `yaml:"custom_commands"`
	AutoBuy        AutoBuy         `yaml:"auto_buy"`
	AutoSell       AutoSell        `yaml:"auto_sell"`
	AutoGift       AutoGift        `yaml:"auto_gift"`
	BalanceCheck   bool            `yaml:"balance_check"`
	LogToFile      bool            `yaml:"log_to_file"`
	Debug          bool            `yaml:"debug"`
}

type AutoGift struct {
	Enable   bool     `yaml:"enable"`
	To       string   `yaml:"to"`
	Interval int      `yaml:"interval"`
	Items    []string `yaml:"items"`
}

type CustomCommand struct {
	Value    string `yaml:"value"`
	Interval int    `yaml:"interval"`
	Amount   int    `yaml:"amount"`
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
	Fish bool `yaml:"fish"`
	Hunt bool `yaml:"hunt"`
}

type SuspicionAvoidance struct {
	Typing       Typing       `yaml:"typing"`
	MessageDelay MessageDelay `yaml:"message_delay"`
}

type Typing struct {
	Base     int `yaml:"base"`     // A base duration in milliseconds.
	Speed    int `yaml:"speed"`    // Speed in keystrokes per minute.
	Variance int `yaml:"variance"` // A random value in milliseconds from [0,n) added to the base.
}

type MessageDelay struct {
	Base     int `yaml:"base"`     // A base duration in milliseconds.
	Variance int `yaml:"variance"` // A random value in milliseconds from [0,n) added to the base.
}

// Shift indicates an application state (active or dormant) for a duration.
type Shift struct {
	State    string   `yaml:"state"`
	Duration Duration `yaml:"duration"`
}

// Duration is not related to a time.Duration. It is a structure used in a Shift
// type.
type Duration struct {
	Base     int `yaml:"base"`     // A base duration in seconds.
	Variance int `yaml:"variance"` // A random value in seconds from [0,n) added to the base.
}

// Load loads the config from the expected path.
func Load(dir string) (Config, error) {
	f, err := os.Open(path.Join(dir, "config.yml"))
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file: %v", err)
	}

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while decoding config: %v", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if len(c.InstancesOpts) == 0 {
		return fmt.Errorf("instances: no instances, at least 1 is required")
	}
	if len(c.Compat.PostmemeOpts) == 0 {
		return fmt.Errorf("compatibility.postmeme: no compatibility options")
	}
	if len(c.Compat.AllowedSearches) == 0 {
		return fmt.Errorf("compatibility.allowed_searches: no compatibility options")
	}
	if len(c.Compat.SearchCancel) == 0 {
		return fmt.Errorf("compatibility.search_cancel: no compatibility options")
	}
	if c.Compat.Cooldown.Postmeme <= 0 {
		return fmt.Errorf("compatibility.cooldown.postmeme: value must be greater than 0")
	}
	if c.Compat.Cooldown.Hunt <= 0 {
		return fmt.Errorf("compatibility.cooldown.hunt: value must be greater than 0")
	}
	if c.Compat.Cooldown.Highlow <= 0 {
		return fmt.Errorf("compatibility.cooldown.highlow: value must be greater than 0")
	}
	if c.Compat.Cooldown.Fish <= 0 {
		return fmt.Errorf("compatibility.cooldown.fish: value must be greater than 0")
	}
	if c.Compat.Cooldown.Search <= 0 {
		return fmt.Errorf("compatibility.cooldown.search: value must be greater than 0")
	}
	if c.Compat.Cooldown.Beg <= 0 {
		return fmt.Errorf("compatibility.cooldown.beg: value must be greater than 0")
	}
	if c.Compat.Cooldown.Margin < 0 {
		return fmt.Errorf("compatibility.cooldown.margin: value must be greater than or equal to 0")
	}
	if c.Compat.AwaitResponseTimeout < 0 {
		return fmt.Errorf("compatibility.await_response_timeout: value must be greater than 0")
	}

	if c.Features.AutoSell.Enable {
		if c.Features.AutoSell.Interval < 0 {
			return fmt.Errorf("features.auto_sell.interval: value must be greater than or equal to 0")
		}
		if len(c.Features.AutoSell.Items) == 0 {
			return fmt.Errorf("features.auto_sell.items: auto_sell enabled but no items configured")
		}
	}

	if c.Features.AutoGift.Enable {
		if c.Features.AutoGift.Interval < 0 {
			return fmt.Errorf("features.auto_gift.interval: value must be greater than or equal to 0")
		}
		if len(c.Features.AutoGift.Items) == 0 {
			return fmt.Errorf("features.auto_gift.items: auto_gift enabled but no items configured")
		}
		if c.Features.AutoGift.To == "" {
			return fmt.Errorf("features.auto_gift.to: no recipient id")
		}
		if !validID(c.Features.AutoGift.To) {
			return fmt.Errorf("features.auto_gift.to: invalid id")
		}
	}

	for i, cmd := range c.Features.CustomCommands {
		if cmd.Value == "" {
			return fmt.Errorf("features.custom_commands[%v].value: no value", i)
		}
		if cmd.Amount < 0 {
			return fmt.Errorf("features.custom_commands[%v].amount: value must be greater than or equal to 0", i)
		}
	}

	for i, instance := range c.InstancesOpts {
		if instance.Token == "" {
			return fmt.Errorf("instances[%v]: no token", i)
		}
		if instance.ChannelID == "" {
			return fmt.Errorf("instances[%v]: no channel id", i)
		}
		if !validID(instance.ChannelID) {
			return fmt.Errorf("instances[%v]: invalid channel id", i)
		}
		if len(instance.Shifts) == 0 {
			return fmt.Errorf("instances[%v]: no shifts", i)
		}
		for j, shift := range instance.Shifts {
			if shift.State != ShiftStateActive && shift.State != ShiftStateDormant {
				return fmt.Errorf("instances[%v].shifts[%v]: invalid shift state: %v", i, j, shift.State)
			}
		}
	}
	return nil
}

func validID(id string) bool {
	return regexp.MustCompile(`^[0-9]+$`).Match([]byte(id))
}
