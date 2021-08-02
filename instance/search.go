// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"time"

	"github.com/dankgrinder/dankgrinder/discord"
)

func (in *Instance) search(msg discord.Message) {
	choices := in.returnButtonLabel(3, msg)
	fmt.Println(choices)
	for _, choice := range choices {
		for _, allowed := range in.Compat.AllowedSearches {
			if choice == allowed {
				index := in.returnButtonIndex(choice, 3, msg)
				time.Sleep(2 * time.Second)
				in.pressButton(index, msg)
				fmt.Println(index)
				return
			}
		}
	}
}
