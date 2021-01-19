package main

import "time"

// commands returns a command pointer slice with all commands that should be
// executed periodically. It contains all commands as configured.
func commands() (cmds []*command) {
	cmds = []*command{
		{run: "pls beg", interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin)},
		{run: "pls pm", interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin)},
		{run: "pls search", interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin)},
		{run: "pls hl", interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin)},
	}
	if cfg.Features.Commands.Fish {
		cmds = append(cmds, &command{
			run: "pls fish",
			interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin),
		})
	}
	if cfg.Features.Commands.Hunt {
		cmds = append(cmds, &command{
			run: "pls hunt",
			interval: sec(cfg.Compat.Cooldown.Beg + cfg.Compat.Cooldown.Margin),
		})
	}
	if cfg.Features.BalanceCheck {
		cmds = append(cmds, &command{
			run: "pls bal",
			interval: time.Minute * 2,
		})
	}
	return
}

// asCommands returns a slice of command pointers used in the auto-selling
// feature. The returned slice contains the commands configured to be enabled.
func asCommands() (cmds []*command) {
	if cfg.Features.AutoSell.Boar {
		cmds = append(cmds, &command{
			run: "pls sell boar max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.Dragon {
		cmds = append(cmds, &command{
			run: "pls sell dragon max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.Duck {
		cmds = append(cmds, &command{
			run: "pls sell duck max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.Rabbit {
		cmds = append(cmds, &command{
			run: "pls sell rabbit max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.Skunk {
		cmds = append(cmds, &command{
			run: "pls sell skunk max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.RareFish {
		cmds = append(cmds, &command{
			run: "pls sell rarefish max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.ExoticFish {
		cmds = append(cmds, &command{
			run: "pls sell exoticfish max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.LegendaryFish {
		cmds = append(cmds, &command{
			run: "pls sell legendaryfish max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	if cfg.Features.AutoSell.Fish {
		cmds = append(cmds, &command{
			run: "pls sell fish max",
			interval: sec(cfg.Features.AutoSell.Interval),
		})
	}
	return
}
