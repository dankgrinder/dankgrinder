// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/dankgrinder/dankgrinder/responder"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/scheduler"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

var cfg config.Config

type logFileHook struct {
	dir string
}

func (lfh logFileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (lfh logFileHook) Fire(e *logrus.Entry) error {
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("dankgrinder-%v.log", date)
	f, err := os.OpenFile(path.Join(lfh.dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := (&logrus.JSONFormatter{}).Format(e)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	ex, err := os.Executable()
	if err != nil {
		logrus.Fatalf("could not find executable path: %v", err)
	}
	cfg, err = config.Load(path.Dir(ex))
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}
	if cfg.Features.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if cfg.Features.LogToFile {
		logrus.AddHook(logFileHook{dir: path.Dir(ex)})
	}
	if err = cfg.Validate(); err != nil {
		logrus.Fatalf("invalid config: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	client, err := discord.NewClient(cfg.Token)
	if err != nil {
		logrus.Fatalf("error while creating client: %v", err)
	}
	logrus.Infof("successful authorization as %v", client.User.Username+"#"+client.User.Discriminator)

	cmds := commands()
	var lastState string
	var rspdr *responder.Responder
	var sdlr *scheduler.Scheduler
	for {
		for i, shift := range cfg.SuspicionAvoidance.Shifts {
			dur := shiftDur(shift)
			logrus.StandardLogger().WithFields(map[string]interface{}{
				"state":    shift.State,
				"duration": dur,
			}).Infof("starting shift %v", i+1)
			if shift.State == lastState {
				time.Sleep(dur)
				continue
			}
			lastState = shift.State

			if shift.State == config.ShiftStateDormant {
				if rspdr != nil {
					rspdr.Close()
				}
				if sdlr != nil {
					sdlr.Close()
				}
				time.Sleep(dur)
				continue
			}
			sdlr = &scheduler.Scheduler{
				Client:       client,
				ChannelID:    cfg.ChannelID,
				Typing:       &cfg.SuspicionAvoidance.Typing,
				MessageDelay: &cfg.SuspicionAvoidance.MessageDelay,
			}
			if err = sdlr.Start(); err != nil {
				logrus.Fatalf("error while starting scheduler: %v", err)
			}
			rspdr = &responder.Responder{
				Sdlr:   sdlr,
				Client: client,
				FatalHandler: func(err error) {
					logrus.Fatalf("responder fatal: %v", err)
				},
				ChannelID:       cfg.ChannelID,
				PostmemeOpts:    cfg.Compat.PostmemeOpts,
				AllowedSearches: cfg.Compat.AllowedSearches,
				BalanceCheck:    cfg.Features.BalanceCheck,
				AutoBuy:         &cfg.Features.AutoBuy,
			}
			if err = rspdr.Start(); err != nil {
				logrus.Fatalf("error while starting responder: %v", err)
			}
			for _, cmd := range cmds {
				sdlr.Schedule(cmd, false)
			}
			time.Sleep(dur)
		}
	}
}
