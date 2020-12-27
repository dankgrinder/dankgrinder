package main

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var cfg = config.MustLoad()
var auth = discord.Authorization{Token: cfg.Token}

// cycleTime is how often a command cycle is triggered, where a command cycle
// is a cycle that goes through all configured commands for the bot.
const cycleTime = time.Second * 4

func sendMessage(content string) {
	err := auth.SendMessage(cfg.ChannelID, content, time.Millisecond*300, nil)
	if err != nil {
		logrus.Errorf("%v", err)
	}
}

func main() {
	fmt.Printf("amount: ")
	var s string
	_, err := fmt.Scanln(&s)
	if err != nil {
		logrus.Fatalf("error while scanning stdin: %v", err)
	}
	amount, err := strconv.Atoi(s)
	if err != nil || amount < 1 {
		logrus.Fatalf("invalid input: must be a positive integer")
	}

	t := time.Tick(cycleTime)
	for i := 0; i < amount; i++ {
		go cycle()
		<-t
	}
}
