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

// cycle is a cycle that goes through all configured commands for the bot.
func cycle() {
	logrus.Infof("starting new cycle")
	cmd("pls beg")
	if cfg.Enable.Fish {
		cmd("pls fish")
	}
	if cfg.Enable.Hunt {
		cmd("pls hunt")
	}
	cmd("pls search")
	cmd("pls pm")
	cmd("pls hl")
	if cfg.BalanceCheck.Enable {
		cmd("pls bal")
	}
}
