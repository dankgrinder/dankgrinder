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
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	if cfg.Token == "" {
		logrus.Fatalf("no authorization token configured")
	}
	if cfg.ChannelID == "" {
		logrus.Fatalf("no channel id configured")
	}

	var err error
	user, err = auth.CurrentUser()
	if err != nil {
		logrus.Fatalf("error while getting user information: %v", err)
	}
	logrus.Infof("successful authorization as %v", user.Username+"#"+user.Discriminator)

	connWS()

	sdlr.schedule <- &command{
		content:  "pls beg",
		interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin),
	}
	if cfg.Features.Commands.Fish {
		sdlr.schedule <- &command{
			content:  "pls fish",
			interval: sec(cfg.Compat.Cooldown.Fish + cfg.Compat.Cooldown.Margin),
		}
	}
	if cfg.Features.Commands.Hunt {
		sdlr.schedule <- &command{
			content:  "pls hunt",
			interval: sec(cfg.Compat.Cooldown.Hunt + cfg.Compat.Cooldown.Margin),
		}
	}
	sdlr.schedule <- &command{
		content:  "pls pm",
		interval: sec(cfg.Compat.Cooldown.Postmeme + cfg.Compat.Cooldown.Margin),
	}
	sdlr.schedule <- &command{
		content:  "pls search",
		interval: sec(cfg.Compat.Cooldown.Search + cfg.Compat.Cooldown.Margin),
	}
	sdlr.schedule <- &command{
		content:  "pls hl",
		interval: sec(cfg.Compat.Cooldown.Highlow + cfg.Compat.Cooldown.Margin),
	}
	if cfg.Features.BalanceCheck {
		sdlr.schedule <- &command{
			content:  "pls bal",
			interval: 2 * time.Minute,
		}
	}

	// The main goroutine is permanently dormant here. I have yet to find a way
	// to not make this happen. This is a temporary solution.
	<-make(chan bool)
}
