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

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func isSameLength(s1 string, s2 string) bool { // Function checks length of strings

	if len(s1) == len(s2) {
		return true
	} else {
		return false
	}
}

func hasSameLetters(s1 string, s2 string) bool { // Function checks the letters of the strings (Scrambled and Unscrambled Words)
	var counter int
	for _, letter1 := range s1 {
		for _, letter2 := range s2 {
			if letter1 == letter2 {
				counter++
				break
			}
		}
	}
	if counter == len(s1) {
		return true
	} else {
		return false
	}
}

func (in *Instance) digEventScramble(msg discord.Message) {
	scramble := exp.digEventScramble.FindStringSubmatch(msg.Content)[1] // Scrambled word is at index 1 of the search query

	var Unscrambled []string // Unscrambled Solution

	for _, word := range in.Compat.AllowedScrambles { // Allowed Scrambles in Compatibility
		if isSameLength(scramble, word) && hasSameLetters(scramble, word) {
			Unscrambled = append(Unscrambled, word)
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value: Unscrambled[0], // Note to self: Type error - Value is string and Unscrambled is []String
				Log:   "responding to dig scramble",
			})
			return
		}
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: in.Compat.DigCancel[rand.Intn(len(in.Compat.DigCancel))],
		Log:   "no allowed dig options provided, responding",
	})
}

func (in *Instance) digEnd(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil {
		return
	}
	if msg.ReferencedMessage.Content != digCmdValue {
		return
	}
	if trigger.Value == digCmdValue &&
		!exp.digEventScramble.MatchString(msg.Content) {
		in.sdlr.Resume()
	}
}
