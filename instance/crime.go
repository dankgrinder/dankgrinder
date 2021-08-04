// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"math/rand"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) crime(msg discord.Message) {
	if in.Compat.CrimeMode == 0 {
		choices := in.returnButtonLabel(3, msg)
		for _, choice := range choices {
			for _, allowed := range in.Compat.AllowedCrimes {
				if choice == allowed {
					index := in.returnButtonIndex(choice, 3, msg)
					in.sdlr.ResumeWithCommand(&scheduler.Command{
						Actionrow: 1,
						Button: index,
						Message: msg,
						Log: "Responding to crime from Allowed options",

					})
				}
			}
		}
	} else if in.Compat.CrimeMode == 1 {
		i := rand.Intn(3)
		in.sdlr.ResumeWithCommand(&scheduler.Command{
			Actionrow: 1,
			Button: i+1,
			Message: msg,
			Log: "Responding to crime randomly",
		})

	} else if in.Compat.CrimeMode == 2 {
		choices := in.returnButtonLabel(3, msg)
		for _, choice := range choices {
			for _, allowed := range in.Compat.AllowedCrimes {
				if choice == allowed {
					index := in.returnButtonIndex(choice, 3, msg)
					in.sdlr.ResumeWithCommand(&scheduler.Command{
						Actionrow: 1,
						Button: index,
						Message: msg,
						Log: "Responding to crime from priority options",
					})
					return
				}
			}
		}
		i := rand.Intn(3)
		in.sdlr.ResumeWithCommand(&scheduler.Command{
			Actionrow: 1,
			Button: i+1,
			Message: msg,
			Log: "Responding to crime randomly, no priority option.",
		})
	}
}
