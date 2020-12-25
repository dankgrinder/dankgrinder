package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

func cmd(cmd string) {
	logrus.WithField("command", cmd).Infof("sending command")
	sendMessage(cmd)
	time.Sleep(time.Second * 4)
}

func cycle() {
	if cfg.Token == "" {
		logrus.Fatalf("no authorization token configured")
	}
	if cfg.ChannelID == "" {
		logrus.Fatalf("no channel id configured")
	}
	if cfg.UserID == "" {
		logrus.Fatalf("no user id configured")
	}

	cmd("pls use candy")
}
