// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/responder"
	"github.com/dankgrinder/dankgrinder/scheduler"
	"github.com/sirupsen/logrus"
)

// instance setup?
type instance struct {
	client    *discord.Client
	channelID string
	cmds      []*scheduler.Command
	sdlr      *scheduler.Scheduler
	rspdr     *responder.Responder
	shifts    []config.Shift
	prevState string
	logger    *logrus.Logger
	fatal     chan error
	wg        *sync.WaitGroup
}

// i think this starts the instances
func (in *instance) sleep(dur time.Duration) {
	select {
	case err := <-in.fatal:
		in.logger.Errorf("%v", err)
		runtime.Goexit()
	case <-time.After(dur):
	}
}

// code always conplaning smh
func (in *instance) start() error {
	if in.client == nil {
		return fmt.Errorf("no client")
	}
	if in.channelID == "" {
		return fmt.Errorf("no channel id")
	}
	if len(in.shifts) == 0 {
		return fmt.Errorf("no shifts")
	}
	if in.wg == nil {
		return fmt.Errorf("no waitgroup")
	}
	in.fatal = make(chan error)
	in.logger = logrus.StandardLogger()
	if len(cfg.InstancesOpts) > 1 {
		in.logger = newInstanceLogger(in.client.User.Username, path.Dir(ex)) //lmao
	}
	in.wg.Add(1)
	go func() {
		defer in.wg.Done()
		for {
			for i, shift := range in.shifts {
				dur := shiftDur(shift)
				in.logger.WithFields(map[string]interface{}{
					"state":    shift.State,
					"duration": dur,
				}).Infof("starting shift %v", i+1)
				if shift.State == in.prevState {
					in.sleep(dur)
					continue
				} //.. somthing to do with shifts, I think
				in.prevState = shift.State
				if shift.State == config.ShiftStateDormant {
					if in.rspdr != nil {
						if err := in.rspdr.Close(); err != nil {
							in.logger.Errorf("error while closing responder: %v", err)
						}
					}
					if in.sdlr != nil {
						if err := in.sdlr.Close(); err != nil {
							in.logger.Errorf("error while closing scheduler: %v", err)
						}
					}
					in.sleep(dur)
					continue
				}
				if err := in.startInterface(); err != nil {
					in.logger.Errorf("instance fatal: %v", err)
					return
				}
				for _, cmd := range in.cmds {
					in.sdlr.Schedule(cmd)
				}
				in.sleep(dur)
			}
		}
	}()
	return nil
}

func (in *instance) startInterface() error {
	in.sdlr = &scheduler.Scheduler{
		Client:             in.client,
		ChannelID:          in.channelID,
		Typing:             &cfg.SuspicionAvoidance.Typing,
		MessageDelay:       &cfg.SuspicionAvoidance.MessageDelay,
		Logger:             in.logger,
		AwaitResumeTimeout: sec(cfg.Compat.AwaitResponseTimeout),
		FatalHandler: func(ferr error) {
			if in.rspdr != nil {
				if err := in.rspdr.Close(); err != nil {
					in.logger.Errorf("error while closing responder: %v", err) //...
				}
			}
			in.fatal <- fmt.Errorf("scheduler fatal: %v", ferr)
		},
	}
	if err := in.sdlr.Start(); err != nil {
		return fmt.Errorf("error while starting scheduler: %v", err)
	}
	in.rspdr = &responder.Responder{
		Sdlr:   in.sdlr,
		Client: in.client,
		FatalHandler: func(ferr error) {
			if err := in.sdlr.Close(); err != nil {
				in.logger.Errorf("error while closing scheduler: %v", err) //why am i doing this my life has come to useless comments, whyyyYY /sssss: my god what have I become its takeing over I cant control......go-........͎.͎.͎.͎.͎.̢̝̪̩̌ͭ͛.̢̟̬̞͕̜̭͑.̱̪̬͍ͤ̓̅ͬ͢.̛͚͎͕͎͍͌ͫ̉.̢̤͎̝̳ͤ w̛̭̻̺̉̀ͭh͔͓̦̩̱͊̎̕a̴̟̟̭͈̻͕͐ͫ̋̀t̤̠͍̺̱ͣ͘ ̪̘̫̬͚͚̻̞̆ͫͤ͗̕ī̞͇̭͟s̢͔̯̘͙͓̃ͩ̿ ͍͖̖̯͔̦̳̐̃̌ͯ͘h̓̄͏̜̭ã͍͓̫̮̯͈͖̈́̎́p̄͏̥̪p̷͔̜͎̦͉̭̹͈̏͋͊ͯẹ̸̱͈̈ͅn̩̩ͯ̾̃͞i̘͖̦͉̅̔̅̍͝n̞̰̩͓̠̜̏̃̀g̞̟̞̪̭ͯ̐ͥ͡.͎͔̘̹̼̻̲͔̑͞.̱͇͓̟̯̲̄́͋̚͡.̦̰͔̞̥͇ͧ͒̀ͅ.̛͓͉̯̖̤͓̟̱́ͫ̅.̧͍̟͛.̢̖̭͙͉͒̆̆_̵͎̈́̌ͅ_͓̼̭͕̔ͤ̔ͫ͡_ͥ̔͂҉̹̲̩͙̞̹_̲̼͈̫̼͖͈̖͑͑́h͖͍̦ͫ͢ḛ̻̿̿̆̕l̗̗̜̭̲ͨ̅̀l̨̋ǒͦͤ T͗ͦ͗͝his is not whattt̖́̚tͬt̢͚ͩẗ̔ͨͥ̕t͚̞͉̼̟ͤ͜t̷̖͓̙͑̈́t̸̬̝̻̔t̢͍̗̖̟̆̽͌ͪͅ eW91IHdpbGwgZmluZCB3aGF0IHlvdSBzZWVrIGluIHRoZSB0cmFkZXMgdmFyaWFudCwgdGFsayB0byB0zJPMlcydzLLMpsydzLpozI/Ng82hzKvMvMyczLvMucy7zZPMq2XMiM2uzJrNn82azLvMls2JIOWwuM2uzYvNos2OzK3Mq82UzJzlt6XNi8yIzJLMgM2gzKbMo8yuzZrNlsy55YegzIvNrs2AzKPNmc2OzYjMmeeJh82nzZ3Nh82HzK/Nmsy7zZrMucyYIMyAzanMgM2gzLrNjcypzKblm57NkcyCzYTMvs2fzJbNmc2ZzK/NlsyszKPNlOWHoM2nzaDMs82ZzLrMqcyd44OozYPNpc2tzLXMnc2ZzKQuzInMuMywzKPMssypzYU=
			}
			in.fatal <- fmt.Errorf("responder fatal for %v: %v", in.client.User.Username, ferr)
		},
		ChannelID:       in.channelID,
		PostmemeOpts:    cfg.Compat.PostmemeOpts,
		AllowedSearches: cfg.Compat.AllowedSearches,
		SearchCancel:    cfg.Compat.SearchCancel,
		BalanceCheck:    cfg.Features.BalanceCheck,
		AutoBuy:         &cfg.Features.AutoBuy,
		AutoGift:        &cfg.Features.AutoGift,
		Logger:          in.logger,
	}
	if err := in.rspdr.Start(); err != nil {
		return fmt.Errorf("error while starting responder: %v", err) //
	}
	return nil
}

//heh.
