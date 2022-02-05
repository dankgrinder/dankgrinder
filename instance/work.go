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
	spaces := exp.workEventSoccer.FindStringSubmatch(msg.Content)[1]
	// q is the position of the goal keeper. Finds q and appropriately
	// selects where to shoot.
	var q int = len(spaces)

	if q <= 6 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    3,
			Message:   msg,
			Log:       "responding to work soccer",
		})
	}
	if q > 6 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    1,
			Message:   msg,
			Log:       "responding to work soccer",
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

// Work Event #6 -> Repeat
func (in *Instance) workEventRepeat(msg discord.Message) {
	in.result3 = exp.workEventRepeat.FindStringSubmatch(msg.Content)[1:]

}

func (in *Instance) workEventRetype2(msg discord.Message) {
	first, second, third, fourth, fifth := in.returnButtonIndex(in.result3[0], 5, msg), in.returnButtonIndex(in.result3[1], 5, msg), in.returnButtonIndex(in.result3[2], 5, msg), in.returnButtonIndex(in.result3[3], 5, msg), in.returnButtonIndex(in.result3[4], 5, msg)
	if first != -1 && second != -1 && third != -1 && fourth != -1 && fifth != -1 {
		//	json_msg, _ := json.Marshal(msg)
		//	fmt.Println(string(json_msg))
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      first,
			Message:     msg,
			Log:         "responding to Work Memory",
			AwaitResume: true,
		})

		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      second,
			Message:     msg,
			Log:         "responding to Work Memory",
			AwaitResume: true,
		})

		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      third,
			Message:     msg,
			Log:         "responding to Work Memory",
			AwaitResume: true,
		})

		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      fourth,
			Message:     msg,
			Log:         "responding to Work Memory",
			AwaitResume: true,
		})

		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      fifth,
			Message:     msg,
			Log:         "responding to Work Memory",
			AwaitResume: true,
		})
	}
}

// Work Event #7 -> Color
func (in *Instance) workEventColor(msg discord.Message) {
	colorObject := exp.workEventColor.FindStringSubmatch(msg.Content)[2:]
	// result is a field of Instance struct of type map.
	// Assigning Key - value pairs to the colors and objects
	in.result = map[string]string{colorObject[1]: colorObject[0], colorObject[3]: colorObject[2], colorObject[5]: colorObject[4]}

}

// Work Event #7 -> Color response
func (in *Instance) workEventColor2(msg discord.Message) {
	// Finding target object
	itemObject := exp.workEventColor2.FindStringSubmatch(msg.Content)[1]
	itemColor := in.result[itemObject]
	index := in.returnButtonIndex(itemColor, 4, msg)
	if index != -1 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    index,
			Message:   msg,
			Log:       "Responding to work Color event",
		})
	}
}

// Work Event #8 -> Emoji
func (in *Instance) workEventEmoji(msg discord.Message) {
	in.result2 = exp.workEventEmoji.FindStringSubmatch(msg.Content)[2]
}

func (in *Instance) workEventEmoji2(msg discord.Message) {
	res := in.returnButtonIndex(in.result2, 5, msg)
	if res != -1 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    res,
			Message:   msg,
			Log:       "Responding to work event emoji",
		})
		return
	}
	res2 := in.returnButtonIndex2(in.result2, 5, msg)
	if res2 != -1 {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 2,
			Button:    res2,
			Message:   msg,
			Log:       "Responding to work event emoji",
		})
		return
	}
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
		exp.workEventRepeat.MatchString(msg.Content) ||
		exp.workEventReverse.MatchString(msg.Content) ||
		exp.workEventColor.MatchString(msg.Content) ||
		exp.workEventColor2.MatchString(msg.Content) ||
		exp.workEventSoccer.MatchString(msg.Content) {
		return
	}
	in.sdlr.Resume()
}
