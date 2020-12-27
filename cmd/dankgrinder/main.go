package main

import (
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	cfg   = config.MustLoad()
	user  discord.User
	auth  = discord.Authorization{Token: cfg.Token}
	sched = startNewScheduler()
)

func main() {
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

	sched.schedule <- command{
		content:  "pls beg",
		interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin),
	}
	if cfg.Features.Commands.Fish {
		sched.schedule <- command{
			content:  "pls fish",
			interval: sec(cfg.Compat.Cooldown.Fish + cfg.Compat.Cooldown.Margin),
		}
	}
	if cfg.Features.Commands.Hunt {
		sched.schedule <- command{
			content:  "pls hunt",
			interval: sec(cfg.Compat.Cooldown.Hunt + cfg.Compat.Cooldown.Margin),
		}
	}
	sched.schedule <- command{
		content:  "pls pm",
		interval: sec(cfg.Compat.Cooldown.Postmeme + cfg.Compat.Cooldown.Margin),
	}
	sched.schedule <- command{
		content:  "pls search",
		interval: sec(cfg.Compat.Cooldown.Search + cfg.Compat.Cooldown.Margin),
	}
	sched.schedule <- command{
		content:  "pls hl",
		interval: sec(cfg.Compat.Cooldown.Highlow + cfg.Compat.Cooldown.Margin),
	}
	if cfg.Features.BalanceCheck {
		sched.schedule <- command{
			content:  "pls bal",
			interval: 2 * time.Minute,
		}
	}

	// The main goroutine is permanently dormant here. I have yet to find a way
	// to not make this happen. This is a temporary solution.
	<-make(chan bool)
}
