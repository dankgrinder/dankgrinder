// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var (
	cfg  = config.MustLoad()
	auth = discord.Authorization{Token: cfg.Token}
	user discord.User
)

// cycleTime is how often a command cycle is triggered, where a command cycle
// is a cycle that goes through all configured commands for the bot.
const cycleTime = time.Second * 4

func sendMessage(content string) {
	err := auth.SendMessage(content, discord.SendMessageOpts{
		ChannelID: cfg.ChannelID,
		Typing:    time.Second,
	})
	if err != nil {
		logrus.Errorf("%v", err)
	}
}

func main() {
	logrus.StandardLogger().SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.StandardLogger().SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	if cfg.Token == "" {
		logrus.Fatalf("no authorization token configured")
	}
	if cfg.ChannelID == "" {
		logrus.Fatalf("no channel id configured")
	}

	var err error
	user, err = auth.CurrentUser()
	if err != nil {
		logrus.Fatalf("error while getting user information: %v", err)
	}
	logrus.Infof("successful authorization as %v", user.Username+"#"+user.Discriminator)

	fmt.Printf("amount: ")
	var s string
	_, err = fmt.Scanln(&s)
	if err != nil {
		logrus.Fatalf("error while scanning stdin: %v", err)
	}

	amount, err := strconv.Atoi(s)
	if err != nil || amount < 1 {
		logrus.Fatalf("invalid input: must be a positive integer")
	}

	t := time.Tick(cycleTime)
	for i := 0; i < amount; i++ {
		logrus.Infof("sending command: pls use candy")
		sendMessage("pls use candy")
		<-t
	}
}
