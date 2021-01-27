// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/shiena/ansicolor"
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
	cfg, err := config.Load(path.Dir(ex))
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}
	if cfg.Features.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err = cfg.Validate(); err != nil {
		logrus.Fatalf("invalid config: %v", err)
	}

	fmt.Printf("amount of candy to sell (or 0 for none): ")
	var s string
	_, err = fmt.Scanln(&s)
	if err != nil {
		logrus.Fatalf("error while scanning stdin: %v", err)
	}
	amount, err := strconv.Atoi(s)
	if err != nil || amount < 0 {
		logrus.Fatalf("invalid input: value must be greater than or equal to 0")
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(cfg.Instances))
	for _, instance := range cfg.Instances {
		instance := instance
		go func() {
			defer wg.Done()
			client, err := discord.NewClient(instance.Token)
			if err != nil {
				logrus.Fatalf("error while creating client: %v", err)
			}

			logrus.Infof("successful authorization as %v", client.User.Username+"#"+client.User.Discriminator)

			for i := 0; i < amount; i++ {
				logrus.Infof("sending command: pls use candy")
				if err = client.SendMessage("pls use candy", discord.SendMessageOpts{
					ChannelID: instance.ChannelID,
					Typing:    time.Second * 1,
				}); err != nil {
					logrus.Errorf("%v", err)
				}
				time.Sleep(time.Second * 3)
			}

			for _, cmd := range cfg.Compat.AutoSell {
				cmd = fmt.Sprintf("pls sell %v max", cmd)
				logrus.Infof("sending command: %v", cmd)
				if err = client.SendMessage(cmd, discord.SendMessageOpts{
					ChannelID: instance.ChannelID,
					Typing:    time.Second * 2,
				}); err != nil {
					logrus.Errorf("%v", err)
				}
				time.Sleep(time.Second * 1)
			}
		}()
	}
	wg.Wait()
}
