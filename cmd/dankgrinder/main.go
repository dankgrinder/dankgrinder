// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

var (
	cfg  = config.MustLoad()
	auth = discord.Authorization{Token: cfg.Token}
	sdlr = startNewScheduler()
	user discord.User
)

func main() {
	logrus.StandardLogger().SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.StandardLogger().SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	if cfg.Token == "" {
		logrus.StandardLogger().Fatalf("no authorization token configured")
	}
	if cfg.ChannelID == "" {
		logrus.StandardLogger().Fatalf("no channel id configured")
	}
	if cfg.Features.AutoSell.Interval < 0 {
		logrus.StandardLogger().Fatalf("auto sell interval must be greater than or equal to 0")
	}
	for _, shift := range cfg.SuspicionAvoidance.Shifts {
		if shift.State != config.ShiftStateActive && shift.State != config.ShiftStateDormant {
			logrus.StandardLogger().Fatalf(
				"invalid shift state: %v, allowed options are %v and %v",
				shift.State,
				config.ShiftStateActive,
				config.ShiftStateDormant,
			)
		}
	}

	rand.Seed(time.Now().UnixNano())

	var err error
	user, err = auth.CurrentUser()
	if err != nil {
		logrus.StandardLogger().Fatalf("error while getting user information: %v", err)
	}
	logrus.StandardLogger().Infof("successful authorization as %v", user.Username+"#"+user.Discriminator)

	// Connect to the websocket. The *discord.WSConn can be discarded as it would
	// only be used for closing, but there is no intention to close.
	_, err = discord.NewWSConn(cfg.Token, discord.WSConnOpts{
		MessageRouter: router(),
		ErrHandler:    errHandler,
		FatalHandler:  fatalHandler,
	})
	if err != nil {
		logrus.StandardLogger().Fatalf("%v", err)
	}
	logrus.StandardLogger().Infof("connected to websocket")

	var cmds, asCmds []*command
	var lastState string
	for {
		for i, shift := range cfg.SuspicionAvoidance.Shifts {
			dur := shiftDur(shift)
			logrus.StandardLogger().WithFields(map[string]interface{}{
				"state": shift.State,
				"duration": dur,
			}).Infof("starting shift %v", i + 1)
			if shift.State == lastState {
				time.Sleep(dur)
				continue
			}
			lastState = shift.State
			if shift.State == config.ShiftStateActive {
				cmds = commands()
				asCmds = asCommands()
				for _, cmd := range cmds {
					sdlr.schedule <- cmd
				}
				for _, cmd := range asCmds {
					sdlr.schedule <- cmd
				}
				time.Sleep(dur)
				continue
			}
			for _, cmd := range cmds {
				cmd.interval = 0
			}
			for _, cmd := range asCmds {
				cmd.interval = 0
			}
			time.Sleep(dur)
		}
	}
}
