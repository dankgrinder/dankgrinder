// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"strconv"
	"strings"
)

func (in *Instance) hl(msg discord.Message) {
	if !exp.hl.MatchString(msg.Embeds[0].Description) {
		return
	}
	nstr := strings.Replace(exp.hl.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	n, err := strconv.Atoi(nstr)
	if err != nil {
		in.Logger.Errorf("error while reading highlow hint: %v", err)
		return
	}
	res := "high"
	if n > 50 {
		res = "low"
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: res,
		Log:   "responding to highlow",
	})
}