// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Work Functionality added by https://github.com/V4NSH4J

package instance

import (
	"math/rand"
	"regexp"
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

// Remove punctuation from strings
func removePunctuation(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_ ]+")
	result := reg.ReplaceAllString(s, "")

	return result
}

// Reverse function (Reverses any string)
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Work Event #1 -> Reversing the String
func (in *Instance) workEventReverse(msg discord.Message) {
	frontward := exp.workEventReverse.FindStringSubmatch(msg.Content)[2]

	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: Reverse(frontward),
		Log:   "responding to working event reverse",
	})
}

// Work Event #2 -> Retyping
func (in *Instance) workEventRetype(msg discord.Message) {
	res := exp.workEventRetype.FindStringSubmatch(msg.Content)[2]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to work retype",
	})
}

// Work Event #3 -> Scramble Solve
func (in *Instance) workEventScramble(msg discord.Message) {
	scramble := exp.workEventScramble.FindStringSubmatch(msg.Content)[2]

	var Unscrambled []string

	for _, word := range in.Compat.AllowedScramblesWork {
		if len(scramble) == len(word) && haveSameChars(scramble, word) {
			Unscrambled = append(Unscrambled, word)
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value: Unscrambled[0],
				Log:   "responding to work scramble",
			})
			return
		}
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value:  in.Compat.WorkCancel[rand.Intn(len(in.Compat.WorkCancel))],
		Log:    "no allowed work options provided, responding",
		Amount: 3,
	})
}

// Work Event #4 -> Soccer
func (in *Instance) workEventSoccer(msg discord.Message) {
	spaces := exp.workEventSoccer.FindStringSubmatch(msg.Content)[2]
	// q is the position of the goal keeper. Finds q and appropriately
	// selects where to shoot.
	var q int = len(spaces)

	if q <= 6 {
		in.sdlr.ResumeWithCommand(&scheduler.Command{
			Value: "right",
			Log:   "responding to work soccer",
		})
	}
	if q > 6 {
		in.sdlr.ResumeWithCommand(&scheduler.Command{
			Value: "left",
			Log:   "responding to work soccer",
		})
	}
}

// Work Event #5 -> Hangman
func (in *Instance) workEventHangman(msg discord.Message) {
	hangman := exp.workEventHangman.FindStringSubmatch(msg.Content)[2]
	ree := regexp.MustCompile(`[a-z, A-Z]{1}( _)+`)
	var pruned string = ree.ReplaceAllString(removePunctuation(hangman), `_`)
	var options []string
	for _, x := range in.Compat.AllowedHangman {
		options = append(options, removePunctuation(x))
	}
	_, s := find(pruned, options)
	if len(s) > 0 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Value: s,
			Log:   "responding to Work Hangman",
		})
		return
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value:  in.Compat.WorkCancel[rand.Intn(len(in.Compat.WorkCancel))],
		Log:    "no allowed hangman, cancelling",
		Amount: 3,
	})
}

// Work Event #6 -> Memory
func (in *Instance) workEventMemory(msg discord.Message) {
	words := exp.workEventMemory.FindStringSubmatch(msg.Content)[2:]
	response := strings.Join(words, " ")
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(response),
		Log:   "responding to Work Memory",
	})
}
func (in *Instance) workEventMemory2(msg discord.Message) {
	words := exp.workEventMemory2.FindStringSubmatch(msg.Content)[2:]
	response := strings.Join(words, " ")
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(response),
		Log:   "responding to Work Memory",
	})
}

// Work Event #7 -> Color
func (in *Instance) workEventColor(msg discord.Message) {
	colorObject := exp.workEventColor.FindStringSubmatch(msg.Content)[2:]
	// result is a field of Instance struct of type map.
	// Assigning Key - value pairs to the colors and objects
	in.result = map[string]string{colorObject[1]: colorObject[0], colorObject[3]: colorObject[2], colorObject[5]: colorObject[4]}
	in.sdlr.AwaitResumeTrigger()
}

// Work Event #7 -> Color response
func (in *Instance) workEventColor2(msg discord.Message) {
	// Finding target object
	itemColor := exp.workEventColor2.FindStringSubmatch(msg.Content)[1]
	var res = in.result[itemColor]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: res,
		Log:   "responding to Work Color",
	})
}

// Rework - Incase of promotion
func (in *Instance) workPromotion(msg discord.Message) {
	in.sdlr.ResumeWithCommandOrPrioritySchedule((&scheduler.Command{
		Value: workCmdValue,
		Log:   "Instance promoted, working",
	}))

}

func (in *Instance) WorkEnd(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || trigger.Value != workCmdValue {
		return
	}
	if exp.workEventScramble.MatchString(msg.Content) ||
		exp.workEventRetype.MatchString(msg.Content) ||
		exp.workEventHangman.MatchString(msg.Content) ||
		exp.workEventMemory.MatchString(msg.Content) ||
		exp.workEventMemory2.MatchString(msg.Content) ||
		exp.workEventReverse.MatchString(msg.Content) ||
		exp.workEventColor.MatchString(msg.Content) ||
		exp.workEventColor2.MatchString(msg.Content) ||
		exp.workEventSoccer.MatchString(msg.Content) {
		return
	}
	in.sdlr.Resume()
}
