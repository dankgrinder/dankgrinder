// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"strconv"
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
)

func (in *Instance) blackjack(msg discord.Message) {
	if !strings.Contains(msg.Embeds[0].Author.Name, in.Client.User.Username) {
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
	in.Logger.Infof("calculated blackjack hand value as: %v", hand)

	dealersUpCard := exp.blackjack.FindStringSubmatch(msg.Embeds[0].Fields[1].Value)[1]
	if dealersUpCard == "J" || dealersUpCard == "Q" || dealersUpCard == "K" {
		dealersUpCard = "10"
	}

	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: in.Features.AutoBlackjack.LogicTable[dealersUpCard][hand],
		Log:   "responding to blackjack",
		AwaitResume: true,
	})
}
