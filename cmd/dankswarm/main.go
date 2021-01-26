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
	"path/filepath"
	"sync"
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
	name := fmt.Sprintf("dankswarm-%v.log", date)
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

func startInstances(instances []config.Instance, logDir string) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(len(cfg.Swarm.Instances))
	for _, instance := range instances {
		instance := instance
		go func() {
			defer wg.Done()
			client, err := discord.NewClient(instance.Token)
			if err != nil {
				logrus.Errorf("error while creating client: %v", err)
				return
			}
			logrus.Infof("successful authorization as %v", client.User.Username+"#"+client.User.Discriminator)
			logger := newInstanceLogger(client.User.Username, logDir)
			shiftLoop(client, instance, logger)
		}()
	}
	return wg
}

func shiftLoop(client *discord.Client, instance config.Instance, logger *logrus.Logger) {
	cmds := commands()
	var rspdr *responder.Responder
	var sdlr *scheduler.Scheduler
	var lastState string
	for {
		for i, shift := range instance.Shifts {
			dur := shiftDur(shift)
			fields := map[string]interface{}{
				"state":    shift.State,
				"duration": dur,
			}
			logrus.WithFields(fields).Infof("starting shift %v for %v", i+1, client.User.Username)
			logger.WithFields(fields).Infof("starting shift %v", i+1)
			if shift.State == lastState {
				time.Sleep(dur)
				continue
			}
			lastState = shift.State

			if shift.State == config.ShiftStateDormant {
				if rspdr != nil {
					if err := rspdr.Close(); err != nil {
						logger.Errorf("error while closing responder: %v", err)
					}
				}
				if sdlr != nil {
					if err := sdlr.Close(); err != nil {
						logger.Errorf("error while closing scheduler: %v", err)
					}
				}
				time.Sleep(dur)
				continue
			}

			// If this is reached, shift state must be active.
			sdlr = &scheduler.Scheduler{
				Client:       client,
				ChannelID:    instance.ChannelID,
				Typing:       &cfg.SuspicionAvoidance.Typing,
				MessageDelay: &cfg.SuspicionAvoidance.MessageDelay,
				Logger:       logger,
			}
			if err := sdlr.Start(); err != nil {
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
			if err := rspdr.Start(); err != nil {
				logrus.Errorf("error while starting responder for %v: %v", client.User.Username, err)
				return
			}
			for _, cmd := range cmds {
				sdlr.Schedule(cmd, false)
			}
			time.Sleep(dur)
		}
	}
}

func newInstanceLogger(username, dir string) *logrus.Logger  {
	logger := logrus.New()
	if cfg.Features.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	if cfg.Features.LogToFile {
		logger = logrus.New()
		logger.SetOutput(fileLogger{
			username: username,
			dir:      dir,
		})
		logger.SetFormatter(&logrus.JSONFormatter{})
		return logger
	}

	// To avoid spamming to stdout if logging to a file is turned off,
	// the level is set to error.
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	return logger
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	ex, err := os.Executable()
	if err != nil {
		logrus.Fatalf("could not find executable path: %v", err)
	}
	ex = filepath.ToSlash(ex)
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
	if len(cfg.Swarm.Instances) == 0 {
		logrus.Fatalf("invalid config: swarm.instances: no instances")
	}
	if len(cfg.Swarm.Instances) == 1 {
		logrus.Warnf("you are using swarm mode with only one instance")
	}

	rand.Seed(time.Now().UnixNano())
	startInstances(cfg.Swarm.Instances, path.Dir(ex)).Wait()
}
