// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"math/rand"
)

func (in *Instance) search(msg discord.Message) {
	choices := exp.search.FindStringSubmatch(msg.Content)[1:]
	for _, choice := range choices {
		for _, allowed := range in.Compat.AllowedSearches {
			if choice == allowed {
				in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
					Value: choice,
					Log:   "responding to search",
				})
				return
			}
		}
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: in.Compat.SearchCancel[rand.Intn(len(in.Compat.SearchCancel))],
		Log:   "no allowed search options provided, responding",
	})
}
