package instance

import (
	"math/rand"
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) guessTheNumber(msg discord.Message) {
	if strings.Contains(msg.Content, "You've got 4 attempts to try and guess my random number between") {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Value:       "10",
			Log:         "Responding to GTN with 10",
			AwaitResume: true,
		})
	}
	if exp.guess.MatchString(msg.Content) {
		if exp.guessHint.MatchString(msg.Content) {
			return
		} else {
			a := exp.guess.FindStringSubmatch(msg.Content)[1:] // returns attempts and hints
			if a[0] == "3" && a[1] == "2" {
				in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
					Value:       "hint",
					Log:         "Asking for GTN hint",
					AwaitResume: true,
				})
			}
			if a[0] == "2" && a[1] == "1" {
				in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
					Value:       "hint",
					Log:         "Asking for GTN hint",
					AwaitResume: true,
				})
			}
		}
	}
	if exp.guessHint.MatchString(msg.Content) {
		a := exp.guessHint.FindStringSubmatch(msg.Content)[1:] // Returns Last number, high or low, attempts and hints
		if a[0] == "10" && a[1] == "high" {
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       "5",
				Log:         "Guessing 5",
				AwaitResume: true,
			})
		}
		if a[0] == "10" && a[1] == "low" {
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       "15",
				Log:         "Guessing 15",
				AwaitResume: true,
			})
		}
		if a[0] == "5" && a[1] == "low" {
			p := rand.Intn(4)
			q := rand.Intn(4)
			possibilities := []string{"6", "7", "8", "9"}
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[p],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[q],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
		}
		if a[0] == "5" && a[1] == "high" {
			p := rand.Intn(4)
			q := rand.Intn(4)
			possibilities := []string{"1", "2", "3", "4"}
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[p],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[q],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
		}
		if a[0] == "15" && a[1] == "high" {
			p := rand.Intn(4)
			q := rand.Intn(4)
			possibilities := []string{"14", "13", "12", "11"}
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[p],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[q],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
		}
		if a[0] == "15" && a[1] == "low" {
			p := rand.Intn(4)
			q := rand.Intn(4)
			possibilities := []string{"16", "17", "18", "19"}
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[p],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value:       possibilities[q],
				Log:         "Guessing 15",
				AwaitResume: true,
			})
		}

	}
}

func (in *Instance) gtnEnd(msg discord.Message) {
	in.sdlr.Resume()
}
