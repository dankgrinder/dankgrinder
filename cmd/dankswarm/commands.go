// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"time"

	"github.com/dankgrinder/dankgrinder/scheduler"
)

// commands returns a command pointer slice with all commands that should be
// executed periodically. It contains all commands as configured.
func commands() (cmds []*scheduler.Command) {
	cmds = []*scheduler.Command{
		{Run: "pls beg", Interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin)},
		{Run: "pls pm", Interval: sec(cfg.Compat.Cooldown.Postmeme + cfg.Compat.Cooldown.Margin)},
		{Run: "pls search", Interval: sec(cfg.Compat.Cooldown.Search + cfg.Compat.Cooldown.Margin)},
		{Run: "pls hl", Interval: sec(cfg.Compat.Cooldown.Highlow + cfg.Compat.Cooldown.Margin)},
	}
	if cfg.Features.Commands.Fish {
		cmds = append(cmds, &scheduler.Command{
			Run:      "pls fish",
			Interval: sec(cfg.Compat.Cooldown.Fish + cfg.Compat.Cooldown.Margin),
		})
	}
	if cfg.Features.Commands.Hunt {
		cmds = append(cmds, &scheduler.Command{
			Run:      "pls hunt",
			Interval: sec(cfg.Compat.Cooldown.Hunt + cfg.Compat.Cooldown.Margin),
		})
	}
	if cfg.Features.BalanceCheck {
		cmds = append(cmds, &scheduler.Command{
			Run:      "pls bal",
			Interval: time.Minute * 2,
		})
	}
	return
}
