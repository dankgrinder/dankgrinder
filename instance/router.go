package instance

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const DMID = "270904126974590976"

var exp = struct {
	search,
	fh,
	hl,
	bal,
	gift,
	shop,
	event *regexp.Regexp
}{
	search: regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(.+)\x60,\s\x60(.+)\x60,\s\x60(.+)\x60`),
	fh:     regexp.MustCompile(`10\sseconds.*\s?([Tt]yping|[Tt]ype)\s\x60(.+)\x60`),
	hl:     regexp.MustCompile(`Your hint is \*\*([0-9]+)\*\*`),
	bal:    regexp.MustCompile(`\*\*Wallet\*\*: \x60?⏣?\s?([0-9,]+)\x60?`),
	event:  regexp.MustCompile(`^(Attack the boss by typing|Type) \x60(.+)\x60`),
	gift:   regexp.MustCompile(`[a-zA-Z\s]* \(([0-9]+) owned\)`),
	shop:   regexp.MustCompile(`pls shop ([a-zA-Z\s]+)`),
}

var numFmt = message.NewPrinter(language.English)

func (in *Instance) fh(msg discord.Message) {
	res := exp.fh.FindStringSubmatch(msg.Content)[2]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to fishing or hunting event",
	})
}

func (in *Instance) pm(_ discord.Message) {
	res := in.Compat.PostmemeOpts[rand.Intn(len(in.Compat.PostmemeOpts))]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: res,
		Log:   "responding to postmeme",
	})
}

func (in *Instance) event(msg discord.Message) {
	res := exp.event.FindStringSubmatch(msg.Content)[2]
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to global event",
	})
}

func (in *Instance) search(msg discord.Message) {
	choices := exp.search.FindStringSubmatch(msg.Content)[1:]
	for _, choice := range choices {
		for _, allowed := range in.Compat.AllowedSearches {
			if choice == allowed {
				in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
					Value: choice,
					Log:   "responding to search",
				})
				return
			}
		}
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: in.Compat.SearchCancel[rand.Intn(len(in.Compat.SearchCancel))],
		Log:   "no allowed search options provided, responding",
	})
}

func (in *Instance) hl(msg discord.Message) {
	if !exp.hl.MatchString(msg.Embeds[0].Description) {
		return
	}
	nstr := strings.Replace(exp.hl.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	n, err := strconv.Atoi(nstr)
	if err != nil {
		in.Logger.Errorf("error while reading highlow hint: %v", err)
		return
	}
	res := "high"
	if n > 50 {
		res = "low"
	}
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: res,
		Log:   "responding to highlow",
	})
}

func (in *Instance) balCheck(msg discord.Message) {
	if !strings.Contains(msg.Embeds[0].Title, in.Client.User.Username) {
		return
	}
	if !exp.bal.Match([]byte(msg.Embeds[0].Description)) {
		return
	}
	balstr := strings.Replace(exp.bal.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	balance, err := strconv.Atoi(balstr)
	if err != nil {
		in.Logger.Errorf("error while reading balance: %v", err)
		return
	}
	in.balance = balance
	in.Logger.Infof(
		"current wallet balance: %v coins",
		numFmt.Sprintf("%d", balance),
	)

	if balance > in.Features.AutoShare.MaximumBalance &&
		in.Features.AutoShare.Enable &&
		in.MasterID != "" &&
		in.Client.User.ID != in.MasterID {
		in.sdlr.Schedule(&scheduler.Command{
			Value: fmt.Sprintf(
				"pls share %v <@%v>",
				balance-in.Features.AutoShare.MinimumBalance,
				in.MasterID,
			),
			Log: "sharing all balance above minimum with master instance",
		})
	}

	if in.startingTime.IsZero() {
		in.initialBalance = balance
		in.startingTime = time.Now()
		return
	}
	inc := balance - in.initialBalance
	per := time.Now().Sub(in.startingTime)
	hourlyInc := int(math.Round(float64(inc) / per.Hours()))
	in.Logger.Infof(
		"average income: %v coins/h",
		numFmt.Sprintf("%d", hourlyInc),
	)
}

func (in *Instance) abLaptop(_ discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy laptop",
		Log:   "no laptop, buying a new one",
	})
}

func (in *Instance) abHuntingRifle(_ discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy rifle",
		Log:   "no hunting rifle, buying a new one",
	})
}

func (in *Instance) abFishingPole(_ discord.Message) {
	in.sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy fishing pole",
		Log:   "no fishing pole, buying a new one",
	})
}

func (in *Instance) abTidepod(_ discord.Message) {
	if !strings.Contains(in.sdlr.AwaitResumeTrigger(), "use tide") {
		return
	}
	in.sdlr.Schedule(&scheduler.Command{
		Value:       "pls buy tidepod",
		Log:         "no tidepod, buy a new one",
	})
	in.sdlr.Schedule(&scheduler.Command{
		Value:       "pls use tidepod",
		Log:         "retrying tidepod usage after last unavailability",
		AwaitResume: true,
	})
}

func (in *Instance) gift(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if !strings.Contains(trigger, "shop") {
		return
	}
	if in.Client.User.ID == in.MasterID {
		in.sdlr.Resume()
		return
	}
	if !exp.gift.Match([]byte(msg.Embeds[0].Title)) || !exp.shop.Match([]byte(trigger)) {
		in.sdlr.Resume()
		return
	}
	amount := exp.gift.FindStringSubmatch(msg.Embeds[0].Title)[1]
	item := exp.shop.FindStringSubmatch(trigger)[1]
	
	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: fmt.Sprintf("pls gift %v %v <@%v>", amount, item, in.MasterID),
		Log:   "gifting items",
	})
}

func (in *Instance) tidepod(msg discord.Message) {
	if !strings.Contains(in.sdlr.AwaitResumeTrigger(), "use tide") {
		return
	}

	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: "y",
		Log:   "accepting tidepod",
	})
}

func (in *Instance) tidepodDeath(msg discord.Message) {
	if in.Features.AutoTidepod.BuyLifesaverOnDeath {
		in.sdlr.Schedule(&scheduler.Command{
			Value: "pls buy lifesaver",
			Log:   "buying lifesaver after death from tidepod",
		})
	}
	in.sdlr.Schedule(&scheduler.Command{
		Value: "pls use tidepod",
		Log:   "retrying tidepod usage after previous death",
		AwaitResume: true,
	})
}

// clean removes all characters except for ASCII characters [32, 126] (basically
// all keys you would find on a US keyboard).
func clean(s string) string {
	var result string
	for _, char := range s {
		if regexp.MustCompile(`[\x20-\x7E]`).MatchString(string(char)) {
			result += string(char)
		}
	}
	return result
}

func (in *Instance) router() *discord.MessageRouter {
	rtr := &discord.MessageRouter{}

	// Fishing and hunting events.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fh).
		Mentions(in.Client.User.ID).
		Handler(in.fh)

	// Postmeme.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("What type of meme do you want to post").
		Mentions(in.Client.User.ID).
		Handler(in.pm)

	// Global events.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.event).
		Handler(in.event)

	// Search.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.search).
		Mentions(in.Client.User.ID).
		Handler(in.search)

	// Highlow.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		RespondsTo(in.Client.User.ID).
		Handler(in.hl)

	// Balance report.
	if in.Features.BalanceCheck {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			Handler(in.balCheck)
	}

	// Auto-buy laptop.
	if in.Features.AutoBuy.Laptop {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("oi you need to buy a laptop in the shop to post memes").
			Mentions(in.Client.User.ID).
			Handler(in.abLaptop)
	}

	// Auto-buy fishing pole.
	if in.Features.AutoBuy.FishingPole {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't have a fishing pole").
			Mentions(in.Client.User.ID).
			Handler(in.abFishingPole)
	}

	// Auto-buy hunting rifle.
	if in.Features.AutoBuy.HuntingRifle {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't have a hunting rifle").
			Mentions(in.Client.User.ID).
			Handler(in.abHuntingRifle)
	}

	// Auto-gift
	if in.Features.AutoGift.Enable &&
		in.MasterID != "" &&
		in.MasterID != in.Client.User.ID {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			Handler(in.gift)
	}

	// Auto-tidepod
	if in.Features.AutoTidepod.Enable {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("There's a high chance you'll injure yourself from the tidepod").
			Handler(in.tidepod)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("Eating a tidepod is just dumb and stupid.").
			Handler(in.tidepodDeath)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't own this item??").
			Handler(in.abTidepod)
	}

	return rtr
}
