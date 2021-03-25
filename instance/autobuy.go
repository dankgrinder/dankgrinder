// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) abLaptop(_ discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy laptop",
		Log:   "no laptop, buying a new one",
	})
}

func (in *Instance) abHuntingRifle(_ discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy rifle",
		Log:   "no hunting rifle, buying a new one",
	})
}

func (in *Instance) abFishingPole(_ discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy fishing pole",
		Log:   "no fishing pole, buying a new one",
	})
}

func (in *Instance) abTidepod(_ discord.Message) {
	if !strings.Contains(in.sdlr.AwaitResumeTrigger(), "use tide") {
		return
	}
	in.sdlr.Schedule(&scheduler.Command{
		Value: "pls buy tidepod",
		Log:   "no tidepod, buying a new one",
	})
	in.sdlr.Schedule(&scheduler.Command{
		Value:       "pls use tidepod",
		Log:         "retrying tidepod usage after last unavailability",
		AwaitResume: true,
	})
}
