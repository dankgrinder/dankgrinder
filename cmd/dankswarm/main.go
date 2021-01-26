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
	"github.com/shiena/ansicolor"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/scheduler"
	"github.com/sirupsen/logrus"
)

var cfg config.Config

type fileLogger struct {
	username string
	dir      string
}

func (fl fileLogger) Write(b []byte) (int, error) {
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("dankswarm-%v-%v.log", fl.username, date)
	f, err := os.OpenFile(path.Join(fl.dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(b)
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
	if err = cfg.Validate(); err != nil {
		logrus.Fatalf("invalid config: %v", err)
	}
	if len(cfg.Swarm.Instances) == 0 {
		logrus.Fatalf("invalid config: swarm.instances: no instances")
	}
	if len(cfg.Swarm.Instances) == 1 {
		logrus.Warnf("you are using swarm mode with only one instance")
	}

	rand.Seed(time.Now().UnixNano())
	for _, instance := range cfg.Swarm.Instances {
		instance := instance
		go func() {
			client, err := discord.NewClient(instance.Token)
			if err != nil {
				logrus.Errorf("error while creating client: %v", err)
				return
			}
			logrus.Infof("successful authorization as %v", client.User.Username+"#"+client.User.Discriminator)

			logger := logrus.New()
			logger.SetLevel(logrus.ErrorLevel)
			logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
			logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
			if cfg.Features.LogToFile {
				logger = logrus.New()
				logger.SetOutput(fileLogger{
					username: client.User.Username,
					dir:      path.Dir(ex),
				})
			}
			if cfg.Features.Debug {
				logger.SetLevel(logrus.DebugLevel)
			}

			cmds := commands()
			var rspdr *responder.Responder
			var sdlr *scheduler.Scheduler
			var lastState string
			for {
				for i, shift := range instance.Shifts {
					dur := shiftDur(shift)
					logrus.WithFields(map[string]interface{}{
						"state":    shift.State,
						"duration": dur,
					}).Infof("starting shift %v for %v", i+1, client.User.Username)
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
						ChannelID:    instance.ChannelID,
						Typing:       &cfg.SuspicionAvoidance.Typing,
						MessageDelay: &cfg.SuspicionAvoidance.MessageDelay,
						Logger:       logger,
					}
					if err = sdlr.Start(); err != nil {
						logrus.Errorf("error while starting scheduler for %v: %v", client.User.Username, err)
						return
					}
					rspdr = &responder.Responder{
						Sdlr:   sdlr,
						Client: client,
						FatalHandler: func(err error) {
							logrus.Errorf("responder fatal for %v: %v", client.User.Username, err)
							logger.Errorf("responder fatal: %v", err)
							sdlr.Close()
						},
						ChannelID:       instance.ChannelID,
						PostmemeOpts:    cfg.Compat.PostmemeOpts,
						AllowedSearches: cfg.Compat.AllowedSearches,
						BalanceCheck:    cfg.Features.BalanceCheck,
						AutoBuy:         &cfg.Features.AutoBuy,
						Logger:          logger,
					}
					if err = rspdr.Start(); err != nil {
						logrus.Errorf("error while starting responder for %v: %v", client.User.Username, err)
						return
					}
					for _, cmd := range cmds {
						sdlr.Schedule(cmd, false)
					}
					time.Sleep(dur)
				}
			}
		}()
	}
	<-make(chan bool)
}
