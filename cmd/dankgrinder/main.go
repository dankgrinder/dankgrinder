package main

import (
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var cfg = config.MustLoad()
var user discord.User
var auth = discord.Authorization{Token: cfg.Token}
var sched scheduler

func sendMessage(content string, abort chan bool) {
	delay := ms(cfg.SuspicionAvoidance.MessageDelay.Base)
	if cfg.SuspicionAvoidance.MessageDelay.Variance > 0 {
		delay += ms(rand.Intn(cfg.SuspicionAvoidance.MessageDelay.Variance)) // TODO move delay calculations to new func.
	}
	tt := typingTime(content)
	logrus.WithFields(map[string]interface{}{
		"delay":  delay.String(),
		"typing": tt.String(),
	}).Infof("sending command: %v", content)
	time.Sleep(delay)

	if err := auth.SendMessage(cfg.ChannelID, content, tt, abort); err != nil {
		logrus.Errorf("%v", err)
	}
}

func main() {
	if cfg.Token == "" {
		logrus.Fatalf("no authorization token configured")
	}
	if cfg.ChannelID == "" {
		logrus.Fatalf("no channel id configured")
	} // TODO: move validation to config package

	var err error
	user, err = auth.CurrentUser()
	if err != nil {
		logrus.Fatalf("error while getting user information: %v", err)
	}
	logrus.Infof("successful authorization as %v", user.Username+"#"+user.Discriminator)

	connWS()
	sched = startNewScheduler()
	sched.scheduleInterval("pls beg", sec(cfg.Compat.Cooldown.Beg+cfg.Compat.Cooldown.Margin))
	if cfg.Features.Commands.Fish {
		sched.scheduleInterval("pls fish", sec(cfg.Compat.Cooldown.Fish+cfg.Compat.Cooldown.Margin))
	}
	if cfg.Features.Commands.Hunt {
		sched.scheduleInterval("pls hunt", sec(cfg.Compat.Cooldown.Hunt+cfg.Compat.Cooldown.Margin))
	}
	sched.scheduleInterval("pls pm", sec(cfg.Compat.Cooldown.Postmeme+cfg.Compat.Cooldown.Margin))
	sched.scheduleInterval("pls search", sec(cfg.Compat.Cooldown.Search+cfg.Compat.Cooldown.Margin))
	sched.scheduleInterval("pls hl", sec(cfg.Compat.Cooldown.Highlow+cfg.Compat.Cooldown.Margin))
	if cfg.Features.BalanceCheck {
		sched.scheduleInterval("pls bal", 5*time.Minute)
	}

	// The main goroutine is permanently dormant here. I have yet to find a way
	// to not make this happen. This is a temporary solution.
	<-make(chan bool)
}
