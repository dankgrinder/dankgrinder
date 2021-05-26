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

func find(tt string, stuff []string) (string, string) { //Some amazing stolen/borrowed code which finds the missing word and phrase from a list of pre-determined phrases
	re := regexp.MustCompile("_")
	q := re.ReplaceAllString(tt, `(\w+)`)
	re2 := regexp.MustCompile(q)
	for _, s := range stuff {
		found := re2.FindStringSubmatch(s)
		if len(found) > 0 {
			return s, found[1]
		}
	}
	return "", ""
}

func (in *Instance) digEventScramble(msg discord.Message) { // Dig Event Unscramble
	scramble := exp.digEventScramble.FindStringSubmatch(msg.Content)[1] // Scrambled word is at index 1 of the search query

	var Unscrambled []string // Unscrambled Solution

	for _, word := range in.Compat.AllowedScrambles { // Allowed Scrambles in Compatibility
		if isSameLength(scramble, word) && hasSameLetters(scramble, word) {
			Unscrambled = append(Unscrambled, word)
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value: Unscrambled[0],
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

func (in *Instance) digEventRetype(msg discord.Message) { // Dig Event Retype Stolen from Grind95
	res := exp.digEventRetype.FindStringSubmatch(msg.Content)[1]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to Dig retype event",
	})
}

func (in *Instance) digEventFTB(msg discord.Message) { // Dig event Fill in the blank
	filltheblank := exp.digEventFTB.FindStringSubmatch(msg.Content)[1]
	ree := regexp.MustCompile(`[a-z]{1}( _)+`)
	var pruned string = ree.ReplaceAllString(filltheblank, `_`) //Regex tested, is correct.
	_, s := find(pruned, in.Compat.AllowedFTB)
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: s,
		Log:   "responding to Dig retype event",
	})
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: in.Compat.DigCancel[rand.Intn(len(in.Compat.DigCancel))],
		Log:   "no allowed fill in the blanks, cancelling",
	})
}

func (in *Instance) digEnd(msg discord.Message) { // Scheduling
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
	if trigger.Value == digCmdValue &&
		!exp.digEventRetype.MatchString(msg.Content) {
		in.sdlr.Resume()
	}
	if trigger.Value == digCmdValue &&
		!exp.digEventFTB.MatchString(msg.Content) {
		in.sdlr.Resume()
	}
}
