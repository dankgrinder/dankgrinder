package instance

import (
	"math/rand"
	"regexp"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const DMID = "270904126974590976"

var exp = struct {
	search,
	fhEvent,
	hl,
	bal,
	gift,
	shop,
	blackjack,
	blackjackBal,
	event *regexp.Regexp
}{
	search:       regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(.+)\x60,\s\x60(.+)\x60,\s\x60(.+)\x60`),
	fhEvent:      regexp.MustCompile(`10\sseconds.*\s?([Tt]yping|[Tt]ype)\s\x60(.+)\x60`),
	hl:           regexp.MustCompile(`Your hint is \*\*([0-9]+)\*\*`),
	bal:          regexp.MustCompile(`\*\*Wallet\*\*: \x60?⏣?\s?([0-9,]+)\x60?`),
	event:        regexp.MustCompile(`^(Attack the boss by typing|Type) \x60(.+)\x60`),
	gift:         regexp.MustCompile(`[a-zA-Z\s]* \(([0-9,]+) owned\)`),
	shop:         regexp.MustCompile(`pls shop ([a-zA-Z\s]+)`),
	blackjack:    regexp.MustCompile(`\x60[♥♦♠♣] ([0-9]{1,2}|[JQKA])\x60`),
	blackjackBal: regexp.MustCompile(`(You now have|You have) (\*\*)?(⏣\s)?(\*\*)?([0-9,]+)(\*\*)?(\sstill)?\.`),
}

var numFmt = message.NewPrinter(language.English)

func (in *Instance) fhEvent(msg discord.Message) {
	res := exp.fhEvent.FindStringSubmatch(msg.Content)[2]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to fishing or hunting event",
	})
}

func (in *Instance) fhEnd(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil {
		return
	}
	if msg.ReferencedMessage.Content != fishCmdValue && msg.ReferencedMessage.Content != huntCmdValue {
		return
	}
	if trigger.Value == fishCmdValue || trigger.Value == huntCmdValue &&
		!exp.fhEvent.MatchString(msg.Content) {
		in.sdlr.Resume()
	}
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
		Log:   "responding to event",
	})
}

// clean removes all characters except for ASCII characters [32, 126] (basically
// all keys you would find on a US keyboard).
func clean(s string) string {
	var result string
	allowedChars := regexp.MustCompile(`[\x20-\x7E]`)
	for _, char := range s {
		if allowedChars.MatchString(string(char)) {
			result += string(char)
		}
	}
	return result
}

func (in *Instance) router() *discord.MessageRouter {
	rtr := &discord.MessageRouter{}

	// Fishing and hunting.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fhEvent).
		Mentions(in.Client.User.ID).
		Handler(in.fhEvent)

	// When a fish/hunt is completed without any events, Dank Memer will
	// reference the original command. If there are events it will mention. This
	// can therefore be used to differentiate between the two.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		RespondsTo(in.Client.User.ID).
		Handler(in.fhEnd)

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
			Handler(in.abFishingPole)
	}

	// Auto-buy hunting rifle.
	if in.Features.AutoBuy.HuntingRifle {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't have a hunting rifle").
			Handler(in.abHuntingRifle)
	}

	// Auto-gift
	if in.Features.AutoGift.Enable &&
		in.Master != nil &&
		in != in.Master {
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
			ContentContains("You lost **all of your coins**.").
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

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			Handler(in.blackjackEnd)
	}

	return rtr
}
