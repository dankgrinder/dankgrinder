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

func (in *Instance) tidepod(_ discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != tidepodCmdValue {
		return
	}

	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: acceptTidepodCmdValue,
		Log:   "accepting tidepod",
	})
}

func (in *Instance) tidepodDeath(_ discord.Message) {
	if in.Features.AutoTidepod.BuyLifesaverOnDeath {
		in.sdlr.Schedule(&scheduler.Command{
			Value: buyCmdValue("1", "lifesaver"),
			Log:   "buying lifesaver after death from tidepod",
		})
	}
	in.sdlr.Schedule(&scheduler.Command{
		Value:       tidepodCmdValue,
		Log:         "retrying tidepod usage after previous death",
		AwaitResume: true,
	})
}
