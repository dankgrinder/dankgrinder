package main

import (
	"github.com/dankgrinder/dankgrinder/api"
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/sirupsen/logrus"
	"time"
)

// cycleTime is how often a command cycle is triggered, where a command cycle
// is a cycle that goes through all configured commands for the bot.
const cycleTime = time.Second * 70

var cfg = config.MustLoad()

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
	if cfg.Token == "" {
		logrus.Fatalf("no authorization token configured")
	}
	if cfg.ChannelID == "" {
		logrus.Fatalf("no channel id configured")
	}
	if cfg.UserID == "" {
		logrus.Fatalf("no user id configured")
	}

	connWS()
	t := time.Tick(cycleTime)
	for i := 0; true; i++ {
		cycle()
		<-t
	}
}
