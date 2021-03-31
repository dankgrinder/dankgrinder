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

func (in *Instance) gift(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || !strings.Contains(trigger.Value, shopBaseCmdValue) {
		return
	}
	if in == in.Master {
		in.sdlr.Resume()
		return
	}
	if !exp.gift.Match([]byte(msg.Embeds[0].Title)) || !exp.shop.Match([]byte(trigger.Value)) {
		in.sdlr.Resume()
		return
	}
	amount := strings.Replace(exp.gift.FindStringSubmatch(msg.Embeds[0].Title)[1], ",", "", -1)
	item := exp.shop.FindStringSubmatch(trigger.Value)[1]

	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: giftCmdValue(amount, item, in.Master.Client.User.ID),
		Log:   "gifting items",
	})
}
