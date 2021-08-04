// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"strconv"
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) blackjack(msg discord.Message) {

	if !strings.Contains(clean(msg.Embeds[0].Author.Name), in.Client.User.Username) {
		return
	}
	if len(msg.Embeds[0].Fields) != 2 {
		return
	}
	var handRaw []string
	for _, match := range exp.blackjack.FindAllStringSubmatch(msg.Embeds[0].Fields[0].Value, -1) {
		if match[1] == "J" || match[1] == "Q" || match[1] == "K" {
			handRaw = append(handRaw, "10")
			continue
		}
		handRaw = append(handRaw, match[1])
	}
	handValue, handAcesAmount, handIsSoft := 0, 0, false
	for _, card := range handRaw {
		if card == "A" {
			handAcesAmount++
			continue
		}
		if card == "J" || card == "Q" || card == "K" {
			handValue += 10
			continue
		}
		cardValue, err := strconv.Atoi(card)
		if err != nil {
			in.Logger.Errorf("error while counting cards: unexpected card: %v", cardValue)
			return
		}
		handValue += cardValue
	}
	for highAcesAmount := handAcesAmount; highAcesAmount >= 0; highAcesAmount-- {
		if handValue+highAcesAmount*11+handAcesAmount-highAcesAmount > 21 {
			if highAcesAmount == 0 {
				in.Logger.Errorf("error while counting aces: hand is bust")
				return
			}
			continue
		}
		handValue += handAcesAmount - highAcesAmount
		handValue += highAcesAmount * 11
		handIsSoft = highAcesAmount != 0
		break
	}

	hand := strconv.Itoa(handValue)
	if handIsSoft {
		hand = "soft" + hand
	}

	dealersUpCard := exp.blackjack.FindStringSubmatch(msg.Embeds[0].Fields[1].Value)[1]
	if dealersUpCard == "J" || dealersUpCard == "Q" || dealersUpCard == "K" {
		dealersUpCard = "10"
	}

	in.Logger.Infof("calculated blackjack hand as: %v against dealer's %v", hand, dealersUpCard)

	if in.Features.AutoBlackjack.LogicTable[dealersUpCard][hand] == "h" {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow:   1,
			Button:      1,
			Message:     msg,
			Log:         "Responding with hit",
			AwaitResume: true,
		})
	}
	if in.Features.AutoBlackjack.LogicTable[dealersUpCard][hand] == "s" {
		in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
			Actionrow: 1,
			Button:    2,
			Message:   msg,
			Log:       "Responding with stand",
		})
	}
}

func (in *Instance) blackjackEnd(msg discord.Message) {
	if strings.Contains(msg.Content, "Type `h` to **hit**, type `s` to **stand**, or type `e` to **end** the game.") {
		return
	}
	if !strings.Contains(clean(msg.Embeds[0].Author.Name), in.Client.User.Username) {
		return
	}
	if !strings.Contains(msg.Embeds[0].Author.Name, "blackjack") {
		return
	}
	if !exp.blackjackBal.MatchString(msg.Embeds[0].Description) {
		return
	}
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger != nil {
	rowLoop:
		for _, row := range in.Features.AutoBlackjack.LogicTable {
			for _, val := range row {
				if val == trigger.Value {
					in.sdlr.Resume()
					break rowLoop
				}
			}
		}
	}
	balstr := strings.Replace(exp.blackjackBal.FindStringSubmatch(msg.Embeds[0].Description)[5], ",", "", -1)
	balance, err := strconv.Atoi(balstr)
	if err != nil {
		in.Logger.Errorf("error while reading balance: %v", err)
		return
	}
	in.updateBalance(balance)
}
