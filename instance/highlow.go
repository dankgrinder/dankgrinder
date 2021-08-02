// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"time"

	"github.com/dankgrinder/dankgrinder/discord"
)

func (in *Instance) hl(msg discord.Message) {
	hint := exp.hl.FindStringSubmatch(msg.Embeds[0].Description)[1]
	if hint[0] > 50 {
		time.Sleep(1 * time.Second)
		in.pressButton(0, msg)
	}
	if hint[0] < 50 {
		time.Sleep(1 * time.Second)
		in.pressButton(2, msg)
	}
	if hint[0] == 50 {
		time.Sleep(1 * time.Second)
		in.pressButton(1, msg)
	}
}
