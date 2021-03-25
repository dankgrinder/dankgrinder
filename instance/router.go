package instance

import (
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math/rand"
	"regexp"
)

const DMID = "270904126974590976"

var exp = struct {
	search,
	fh,
	hl,
	bal,
	gift,
	shop,
	blackjack,
	event *regexp.Regexp
}{
	search: regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(.+)\x60,\s\x60(.+)\x60,\s\x60(.+)\x60`),
	fh:     regexp.MustCompile(`10\sseconds.*\s?([Tt]yping|[Tt]ype)\s\x60(.+)\x60`),
	hl:     regexp.MustCompile(`Your hint is \*\*([0-9]+)\*\*`),
	bal:    regexp.MustCompile(`\*\*Wallet\*\*: \x60?⏣?\s?([0-9,]+)\x60?`),
	event:  regexp.MustCompile(`^(Attack the boss by typing|Type) \x60(.+)\x60`),
	gift:   regexp.MustCompile(`[a-zA-Z\s]* \(([0-9,]+) owned\)`),
	shop:   regexp.MustCompile(`pls shop ([a-zA-Z\s]+)`),
	blackjack: regexp.MustCompile(`\[\x60.\s([0-9]{1,2}|[JQKA])\x60\]`),
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
		HasEmbeds(false).
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
	if in.Features.BalanceCheck.Enable {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			Handler(in.balanceCheck)
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

	// Auto-blackjack
	if in.Features.AutoBlackjack.Enable {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			ContentContains("Type `h` to **hit**, type `s` to **stand**, or type `e` to **end** the game.").
			Handler(in.blackjack)
	}

	return rtr
}
