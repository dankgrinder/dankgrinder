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

func haveSameChars(s1 string, s2 string) bool {
	dissect1, dissect2 := map[rune]int{}, map[rune]int{}
	for _, c := range s1 {
		dissect1[c]++
	}
	for _, c := range s2 {
		dissect2[c]++
	}
	if len(dissect1) != len(dissect2) {
		return false
	}
	for c, n := range dissect1 {
		if dissect2[c] != n {
			return false
		}
	}
	return true
}

// Function compares an input phrase with a missing word
// against a list of complete phrases and singles out the
// missing word.
func find(tt string, listOfPhrases []string) (string, string) {
	re := regexp.MustCompile("_")
	q := re.ReplaceAllString(tt, `(\w+)`)
	// Replacing the blank with a word character
	re2 := regexp.MustCompile(q)

	// Finding String submatch in the configurable
	// sentences and returning the missing word
	for _, s := range listOfPhrases {
		found := re2.FindStringSubmatch(s)
		if len(found) > 0 {
			return s, found[1]
		}
	}
	return "", ""
}

func (in *Instance) digEventScramble(msg discord.Message) {
	scramble := exp.digEventScramble.FindStringSubmatch(msg.Content)[1]

	for _, word := range in.Compat.AllowedScrambles {
		if len(scramble) == len(word) && haveSameChars(scramble, word) {
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value: word,
				Log:   "responding to dig scramble",
			})
			return
		}
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value:  in.Compat.DigCancel[rand.Intn(len(in.Compat.DigCancel))],
		Log:    "no allowed dig options provided, responding",
		Amount: 3,
	})
}

func (in *Instance) digEventRetype(msg discord.Message) {
	res := exp.digEventRetype.FindStringSubmatch(msg.Content)[1]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to Dig retype event",
	})
}

func (in *Instance) digEventFTB(msg discord.Message) {
	fillTheBlank := exp.digEventFTB.FindStringSubmatch(msg.Content)[1]
	ree := regexp.MustCompile(`[a-z, A-Z]{1}( _)+`)
	// Replacing the missing word and the hint with an underscore for compatibility with find function
	var pruned string = ree.ReplaceAllString(fillTheBlank, `_`)
	_, s := find(pruned, in.Compat.AllowedFTB)
	if len(s) > 0 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Value: s,
			Log:   "responding to Dig retype event",
		})
		return
	}

	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value:  in.Compat.DigCancel[rand.Intn(len(in.Compat.DigCancel))],
		Log:    "no allowed fill in the blanks, cancelling",
		Amount: 3,
	})
}

func (in *Instance) digEnd(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != digCmdValue {
		return
	}
	if exp.digEventScramble.MatchString(msg.Content) ||
		exp.digEventRetype.MatchString(msg.Content) ||
		exp.digEventFTB.MatchString(msg.Content) {
		return
	}
	in.sdlr.Resume()
}
