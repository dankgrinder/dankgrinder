// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"strings"
)

func (in *Instance) gift(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if !strings.Contains(trigger, "shop") {
		return
	}
	if in.Client.User.ID == in.MasterID {
		in.sdlr.Resume()
		return
	}
	if !exp.gift.Match([]byte(msg.Embeds[0].Title)) || !exp.shop.Match([]byte(trigger)) {
		in.sdlr.Resume()
		return
	}
	amount := strings.Replace(exp.gift.FindStringSubmatch(msg.Embeds[0].Title)[1], ",", "", -1)
	item := exp.shop.FindStringSubmatch(trigger)[1]

	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: fmt.Sprintf("pls gift %v %v <@%v>", amount, item, in.MasterID),
		Log:   "gifting items",
	})
}