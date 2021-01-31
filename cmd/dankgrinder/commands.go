// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"time"

	"github.com/dankgrinder/dankgrinder/scheduler"
)

// commands returns a command pointer slice with all commands that should be
// executed periodically. It contains all commands as configured.
func commands() (cmds []*scheduler.Command) {
	cmds = []*scheduler.Command{
		{
			Value: "pls beg",
			Interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin),
		},
		{
			Value: "pls pm",
			Interval: sec(cfg.Compat.Cooldown.Postmeme + cfg.Compat.Cooldown.Margin),
			AwaitResume: true,
		},
		{
			Value: "pls search",
			Interval: sec(cfg.Compat.Cooldown.Search + cfg.Compat.Cooldown.Margin),
			AwaitResume: true,
		},
		{
			Value: "pls hl",
			Interval: sec(cfg.Compat.Cooldown.Highlow + cfg.Compat.Cooldown.Margin),
			AwaitResume: true,
		},
	}
	if cfg.Features.Commands.Fish {
		cmds = append(cmds, &scheduler.Command{
			Value:    "pls fish",
			Interval: sec(cfg.Compat.Cooldown.Fish + cfg.Compat.Cooldown.Margin),
			AwaitResume: true,
		})
	}
	if cfg.Features.Commands.Hunt {
		cmds = append(cmds, &scheduler.Command{
			Value:    "pls hunt",
			Interval: sec(cfg.Compat.Cooldown.Hunt + cfg.Compat.Cooldown.Margin),
			AwaitResume: true,
		})
	}
	if cfg.Features.BalanceCheck {
		cmds = append(cmds, &scheduler.Command{
			Value:    "pls bal",
			Interval: time.Minute * 2,
		})
	}
	if cfg.Features.AutoSell.Enable {
		var sellCmds []*scheduler.Command
		for i, item := range cfg.Features.AutoSell.Items {
			sellCmds = append(sellCmds, &scheduler.Command{
				Value: fmt.Sprintf("pls sell %v max", item),
				Interval: time.Second * 5,
			})
			if i != 0 {
				sellCmds[i - 1].Next = sellCmds[i]
			}
			if i == len(cfg.Features.AutoSell.Items) - 1 {
				sellCmds[i].Next = sellCmds[0]
				sellCmds[i].Interval = sec(cfg.Features.AutoSell.Interval)
			}
		}
		cmds = append(cmds, sellCmds[0])
	}
	return
}
