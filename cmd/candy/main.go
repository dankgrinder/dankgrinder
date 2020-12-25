package main

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/api"
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var cfg = config.MustLoad()

// cycleTime is how often a command cycle is triggered, where a command cycle
// is a cycle that goes through all configured commands for the bot.
const cycleTime = time.Second * 4

func sendMessage(content string) {
	if err := api.SendMessage(api.SendMessageOpts{
		Token:     cfg.Token,
		ChannelID: cfg.ChannelID,
		Content:   content,
		Typing:    time.Duration(cfg.TypingDuration) * time.Millisecond,
	}); err != nil {
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
