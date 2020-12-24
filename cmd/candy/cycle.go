package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

func cmd(cmd string) {
	logrus.WithField("command", cmd).Infof("sending command")
	sendMessage(cmd)
	time.Sleep(time.Duration(cfg.CmdDelay) * time.Millisecond)
}

func cycle() {
	cmd("pls use candy")
}