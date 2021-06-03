// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Dig Functionality added by https://github.com/V4NSH4J

package instance

import (
	"math/rand"
	"regexp"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) fishEventScramble(msg discord.Message) { //
	scramble := exp.fishEventScramble.FindStringSubmatch(msg.Content)[1]

	for _, word := range in.Compat.AllowedScramblesFish { //
		if len(scramble) == len(word) && haveSameChars(scramble, word) {
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value: word,
				Log:   "responding to fish scramble",
			})
			return
		}
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value:  in.Compat.FishCancel[rand.Intn(len(in.Compat.FishCancel))], //
		Log:    "no allowed fish options provided, responding",
		Amount: 3,
	})
}

func (in *Instance) fishEventFTB(msg discord.Message) { //
	fillTheBlank := exp.fishEventFTB.FindStringSubmatch(msg.Content)[1]
	ree := regexp.MustCompile(`[a-z, A-Z]{1}( _)+`)
	// Replacing the missing word and the hint with an underscore for compatibility with find function
	var pruned string = ree.ReplaceAllString(fillTheBlank, `_`)
	_, s := find(pruned, in.Compat.AllowedFishFTB) //
	if len(s) > 0 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Value: s,
			Log:   "responding to Fish fill the blank",
		})
		return
	}

	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value:  in.Compat.FishCancel[rand.Intn(len(in.Compat.FishCancel))],
		Log:    "no allowed fill in the blanks, cancelling",
		Amount: 3,
	})
}

func (in *Instance) fishEventReverse(msg discord.Message) {
	frontward := exp.fishEventReverse.FindStringSubmatch(msg.Content)[1]

	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: Reverse(frontward),
		Log:   "responding to fish event reverse",
	})
}

func (in *Instance) fishEventRetype(msg discord.Message) {
	res := exp.fishEventRetype.FindStringSubmatch(msg.Content)[1]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to fishing event",
	})
}

func (in *Instance) fishEnd(msg discord.Message) { //
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != fishCmdValue {
		return
	}
	if exp.fishEventScramble.MatchString(msg.Content) ||
		exp.fishEventReverse.MatchString(msg.Content) ||
		exp.fishEventRetype.MatchString(msg.Content) ||
		exp.fishEventFTB.MatchString(msg.Content) {
		return
	}
	in.sdlr.Resume()
}
