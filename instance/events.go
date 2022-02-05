package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) huntEvent(msg discord.Message) {
	fmt.Println("Animal")
	fireball := exp.huntEvent.FindStringSubmatch(msg.Content)[2]
	levitate := exp.huntEvent.FindStringSubmatch(msg.Content)[1]
	posLevitate := len(levitate)
	posFireball := len(fireball)
	if posFireball > posLevitate {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      1,
			Message:     msg,
			Log:         "Catching animal",
			AwaitResume: true,
		})
	}
	if posFireball < posLevitate {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      3,
			Message:     msg,
			Log:         "Catching animal",
			AwaitResume: true,
		})
	}
	if posFireball == posLevitate {
		i := rand.Intn(1)
		if i == 0 {
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Actionrow:   1,
				Button:      1,
				Message:     msg,
				Log:         "Catching animal",
				AwaitResume: true,
			})
		} else {
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Actionrow:   1,
				Button:      3,
				Message:     msg,
				Log:         "Catching animal",
				AwaitResume: true,
			})
		}
	}

}

func (in *Instance) huntEnd(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil {
		return
	}
	if msg.ReferencedMessage.Content != huntCmdValue {
		return
	}
	if trigger.Value == huntCmdValue &&
		!exp.huntEvent.MatchString(msg.Content) {
		in.sdlr.Resume()
	}
}

func (in *Instance) event(msg discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Actionrow: 1,
		Button:    1,
		Message:   msg,
		Log:       "responding to event",
	})
}

func (in *Instance) shopEvent(msg discord.Message) {
	eventType := exp.shopEvent.FindStringSubmatch(msg.Embeds[0].Description)[1]
	json_msg, _ := json.Marshal(msg)
	fmt.Println(string(json_msg))
	type ShopItems struct {
		ID      string `json:"ID"`
		EmojiID string `json:"EmojiID"`
		Name    string `json:"Name"`
		Cost    string `json:"Cost"`
		Type    string `json:"Type"`
	}
	type Shop struct {
		Shop []ShopItems `json:"Shop"`
	}
	ex, _ := os.Executable()
	ex = filepath.ToSlash(ex)
	jsonFile, err := os.Open((path.Join(path.Dir(ex), "shop.json")))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	bytevalue, _ := ioutil.ReadAll(jsonFile)

	var shop Shop
	json.Unmarshal(bytevalue, &shop)
	url := regexp.MustCompile(`https:\/\/cdn.discordapp.com\/emojis\/(.+)\.png`)
	if url.Match([]byte(msg.Embeds[0].Image.URL)) {
		picture := url.FindStringSubmatch(msg.Embeds[0].Image.URL)[1]

		if eventType == "name" {
			for i := 0; i < len(shop.Shop); i++ {
				if picture == shop.Shop[i].EmojiID {
					name := shop.Shop[i].Name
					index := in.returnButtonIndex(name, 4, msg)
					if index != -1 {
						in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
							Actionrow: 1,
							Button:    index,
							Message:   msg,
							Log:       "Responding to shop global event",
						})
					}
				}

			}
		}

		if eventType == "cost" {
			for i := 0; i < len(shop.Shop); i++ {
				if picture == shop.Shop[i].EmojiID {
					cost := shop.Shop[i].Cost
					index := in.returnButtonIndex(cost, 4, msg)
					if index != -1 {
						in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
							Actionrow: 1,
							Button:    index,
							Message:   msg,
							Log:       "Responding to shop global event",
						})
					}
				}
			}
		}
		if eventType == "type" {
			for i := 0; i < len(shop.Shop); i++ {
				if picture == shop.Shop[i].EmojiID {
					typeItem := shop.Shop[i].Type
					index := in.returnButtonIndex(typeItem, 4, msg)
					if index != -1 {
						in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
							Actionrow: 1,
							Button:    index,
							Message:   msg,
							Log:       "Responding to shop global event",
						})
					}
				}
			}
		}
	}
}

func (in *Instance) fishCatch(msg discord.Message) {
	fmt.Println("fish")
	if exp.fishCatch2.Match([]byte(msg.Content)) {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    1,
			Message:   msg,
			Log:       "Catching fish",
		})
		return
	}
	spaces := exp.fishCatch.FindStringSubmatch(msg.Content)[1]
	position := len(spaces)
	if position == 8 && position < 13 {
		in.sdlr.PrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    2,
			Message:   msg,
			Log:       "Catching fish",
		})
	}
	if position == 15 {
		in.sdlr.PrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    3,
			Message:   msg,
			Log:       "Catching fish",
		})
	}
}
