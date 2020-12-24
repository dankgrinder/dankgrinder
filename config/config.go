package config

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

const fileName = "config.json"

type (
	Config struct {
		Token          string       `json:"token"`
		ChannelID      string       `json:"channel_id"`
		UserID         string       `json:"user_id"`
		ResDelay       int          `json:"response_delay"`
		CmdDelay       int          `json:"command_delay"`
		TypingDuration int          `json:"typing_duration"`
		Enable         Enable       `json:"enable"`
		Postmeme       []string     `json:"postmeme"`
		GlobalEvents   []string     `json:"global_events"`
		Search         []string     `json:"search"`
		BalanceCheck   BalanceCheck `json:"balance_check"`
		AutoBuy        AutoBuy      `json:"auto_buy"`
	}
	Enable struct {
		Fish bool `json:"fish"`
		Hunt bool `json:"hunt"`
	}
	BalanceCheck struct {
		Enable   bool   `json:"enable"`
		Username string `json:"username"`
	}
	AutoBuy struct {
		FishingPole  bool `json:"fishing_pole"`
		HuntingRifle bool `json:"hunting_rifle"`
		Laptop       bool `json:"laptop"`
	}
)

func configPath() string {
	ex, err := os.Executable()
	if err != nil {
		logrus.Fatalf("cannot find executable path: %v", err)
	}
	return path.Dir(ex)
}

// MustLoad runs Load and calls logrus.Fatalf is an error occurs.
func MustLoad() Config {
	c, err := Load()
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	return c
}

// Load loads the config from the expected path.
func Load() (Config, error) {
	f, err := os.Open(path.Join(configPath(), fileName))
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file: %v", err)
	}
	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("error while decoding config as json: %v", err)
	}
	return c, nil
}
