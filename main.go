// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/dankgrinder/dankgrinder/instance"

	"github.com/shiena/ansicolor"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	ex, err := os.Executable()
	if err != nil {
		logrus.Fatalf("could not find executable path: %v", err)
	}
	ex = filepath.ToSlash(ex)

	var cfg config.Config
	if len(os.Args) > 1 {
		logrus.Infof("loading config from %v", os.Args[1])
		cfg, err = config.Load(os.Args[1])
	} else {
		logrus.Infof("loading config from %v", path.Join(path.Dir(ex), "config.yml"))
		cfg, err = config.Load(path.Join(path.Dir(ex), "config.yml"))
	}
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}
	if cfg.Features.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.AddHook(logFileHook{dir: path.Dir(ex)})

	// Checks for many possible invalid configurations. This means that during
	// execution of the program, many of these checks don't need to be repeated.
	if err = cfg.Validate(); err != nil {
		logrus.Fatalf("invalid config: %v", err)
	}
	if cfg.Compat.AwaitResponseTimeout < 3 {
		logrus.Warnf("await response timeout is less than 3, this might cause stability issues for responses")
	}

	rand.Seed(time.Now().UnixNano())

	wg := &sync.WaitGroup{}
	for ck, cluster := range cfg.Clusters {
		var ins []*instance.Instance
		var master *instance.Instance
		logrus.Infof("starting cluster %v", ck)

		for i, inOpts := range append(cluster.Instances, cluster.Master) {
			client, err := discord.NewClient(inOpts.Token)
			if err != nil {
				logrus.Errorf("error while creating client: %v", err)
				if i == 0 {
					logrus.Warnf("failed to create master instance client, some functionality may be unavailable")
				}
				continue
			}

			logrus.Infof("successful authorization as %v", client.User.Username+"#"+client.User.Discriminator)

			in := &instance.Instance{
				Client:             client,
				ChannelID:          inOpts.ChannelID,
				WG:                 wg,
				Features:           inOpts.Features,
				SuspicionAvoidance: inOpts.SuspicionAvoidance,
				Compat:             cfg.Compat,
				Shifts:             inOpts.Shifts,
			}

			loggerOpts := instanceLoggerOpts{
				username:             in.Client.User.Username,
				discriminator:        in.Client.User.Discriminator,
				cluster:              ck,
				id:                   in.Client.User.ID,
				debug:                in.Features.Debug,
				verboseStdLoggerHook: in.Features.VerboseLogToStdout,
			}
			if in.Features.LogToFile {
				loggerOpts.dir = path.Dir(ex)
			}
			in.Logger = newInstanceLogger(loggerOpts)

			if i == len(cluster.Instances) {
				master = in
			}
			ins = append(ins, in)
		}

		for _, in := range ins {
			in.Master = master
			in.Cluster = ins
			if err = in.Start(); err != nil {
				logrus.Fatalf("error while starting instance: %v", err)
			}
		}
	}

	wg.Wait()
	logrus.Fatalf("no running instances left")
}
