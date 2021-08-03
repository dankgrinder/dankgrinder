// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"time"
	"math/rand"

	"github.com/dankgrinder/dankgrinder/discord"
)

func (in *Instance) crime(msg discord.Message) {
	if in.Compat.CrimeMode == 0 {
		choices := in.returnButtonLabel(3, msg)
		for _, choice := range choices {
			for _, allowed := range in.Compat.AllowedCrimes {
				if choice == allowed {
					index := in.returnButtonIndex(choice, 3, msg)
					time.Sleep(2 * time.Second)
					in.pressButton(index, msg)
				}
			}
		}
	} else if in.Compat.CrimeMode == 1 {
		i := rand.Intn(3)
		time.Sleep(2 * time.Second)
		in.pressButton(i, msg)

	} else if in.Compat.CrimeMode == 2 {
		choices := in.returnButtonLabel(3, msg)
		for _, choice := range choices {
			for _, allowed := range in.Compat.AllowedCrimes {
				if choice == allowed {
					index := in.returnButtonIndex(choice, 3, msg)
					time.Sleep(2 * time.Second)
					in.pressButton(index, msg)
					return
				}
			}
		}
		i := rand.Intn(3)
		time.Sleep(2 * time.Second)
		in.pressButton(i, msg)
	}
}
