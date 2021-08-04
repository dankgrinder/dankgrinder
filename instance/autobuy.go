// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) abLaptop(_ discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != postmemeCmdValue {
		return
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: buyCmdValue("1", "laptop"),
		Log:   "no laptop, buying a new one",
	})
}

func (in *Instance) abShovel(_ discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != digCmdValue {
		return
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: buyCmdValue("1", "shovel"),
		Log:   "no shovel, buying a new one",
	})
}

func (in *Instance) abHuntingRifle(_ discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != huntCmdValue {
		return
	}
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: buyCmdValue("1", "rifle"),
		Log:   "no hunting rifle, buying a new one",
	})
}

func (in *Instance) abFishingPole(_ discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != fishCmdValue {
		return
	}
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: buyCmdValue("1", "fishing"),
		Log:   "no fishing pole, buying a new one",
	})
}

func (in *Instance) abTidepod(_ discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != tidepodCmdValue {
		return
	}
	in.sdlr.Schedule(&scheduler.Command{
		Value: buyCmdValue("1", "tide"),
		Log:   "no tidepod, buying a new one",
	})
	in.sdlr.Schedule(&scheduler.Command{
		Value:       tidepodCmdValue,
		Log:         "retrying tidepod usage after last unavailability",
		AwaitResume: true,
	})
}
