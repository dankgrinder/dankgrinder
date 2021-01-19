// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

const (
	ShiftStateActive  = "active"
	ShiftStateDormant = "dormant"
)

type Config struct {
	Token              string             `yaml:"token"`
	ChannelID          string             `yaml:"channel_id"`
	Features           Features           `yaml:"features"`
	Compat             Compat             `yaml:"compatibility"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
}

type Compat struct {
	Postmeme        []string `yaml:"postmeme"`
	AllowedSearches []string `yaml:"allowed_searches"`
	Cooldown        Cooldown `yaml:"cooldown"`
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
	Commands     Commands `yaml:"commands"`
	AutoBuy      AutoBuy  `yaml:"auto_buy"`
	AutoSell     AutoSell `yaml:"auto_sell"`
	BalanceCheck bool     `yaml:"balance_check"`
}

type AutoBuy struct {
	FishingPole  bool `yaml:"fishing_pole"`
	HuntingRifle bool `yaml:"hunting_rifle"`
	Laptop       bool `yaml:"laptop"`
}

type AutoSell struct {
	Boar          bool `yaml:"boar"`
	Dragon        bool `yaml:"dragon"`
	Duck          bool `yaml:"duck"`
	Fish          bool `yaml:"fish"`
	ExoticFish    bool `yaml:"exotic_fish"`
	LegendaryFish bool `yaml:"legendary_fish"`
	Rabbit        bool `yaml:"rabbit"`
	RareFish      bool `yaml:"rare_fish"`
	Skunk         bool `yaml:"skunk"`
	Interval      int  `yaml:"interval"`
}

type Commands struct {
	Fish bool `yaml:"fish"`
	Hunt bool `yaml:"hunt"`
}

type SuspicionAvoidance struct {
	Typing       Typing       `yaml:"typing"`
	MessageDelay MessageDelay `yaml:"message_delay"`
	Shifts       []Shift      `yaml:"shifts"`
}

type Typing struct {
	Base     int `yaml:"base"`     // A base duration in milliseconds.
	Speed    int `yaml:"speed"`    // Speed in keystrokes per minute.
	Variance int `yaml:"variance"` // A random value in milliseconds from [0,n) added to the base.
}

// MessageDelay is used to
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

func configDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot find executable path: %v", err)
	}
	return path.Dir(ex), nil
}

// MustLoad runs Load and calls logrus.Fatalf is an error occurs.
func MustLoad() Config {
	c, err := Load()
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}
	return c
}

// Load loads the config from the expected path.
func Load() (Config, error) {
	dir, err := configDir()
	if err != nil {
		return Config{}, fmt.Errorf("error while getting config dir: %v", err)
	}
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
