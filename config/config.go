package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

type Config struct {
	Token              string             `yaml:"token"`
	ChannelID          string             `yaml:"channel_id"`
	Features           Features           `yaml:"features"`
	Compat             Compat             `yaml:"compatibility"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
}

type Compat struct {
	Postmeme []string `yaml:"postmeme"`
	Search   []string `yaml:"search"`
	Cooldown Cooldown `yaml:"cooldown"`
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
	BalanceCheck bool     `yaml:"balance_check"`
}

type AutoBuy struct {
	FishingPole  bool `yaml:"fishing_pole"`
	HuntingRifle bool `yaml:"hunting_rifle"`
	Laptop       bool `yaml:"laptop"`
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
	Base     int `yaml:"base"`
	Speed    int `yaml:"speed"`
	Variance int `yaml:"variance"`
}

type MessageDelay struct {
	Base     int `yaml:"base"`
	Variance int `yaml:"variance"`
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
		logrus.Fatalf("%v", err)
	}
	return c
}

// Load loads the config from the expected path.
func Load() (Config, error) {
	dir, err := configDir()
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(path.Join(dir, "config.yml"))
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file: %v", err)
	}

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while unmarshalling config: %v", err)
	}

	return cfg, nil
}
