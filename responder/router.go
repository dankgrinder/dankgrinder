package responder

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/scheduler"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const DMID = "270904126974590976"

var exp = struct {
	search,
	fh,
	hl,
	bal,
	event *regexp.Regexp
}{
	search: regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(?P<m1>.+)\x60,\s\x60(?P<m2>.+)\x60,\s\x60(?P<m3>.+)\x60`),
	fh:     regexp.MustCompile(`10\sseconds.*\s?([Tt]yping|[Tt]ype)\s\x60(?P<m1>.+)\x60`),
	hl:     regexp.MustCompile(`Your hint is \*\*(?P<m1>[0-9]+)\*\*`),
	bal:    regexp.MustCompile(`\*\*Wallet\*\*: \x60?⏣?\s?(?P<m1>[0-9,]+)\x60?`),
	event:  regexp.MustCompile(`^(Attack the boss by typing|Type) \x60(?P<m1>.+)\x60`),
}

var numFmt = message.NewPrinter(language.English)

func (r *Responder) fh(msg discord.Message) {
	res := exp.fh.FindStringSubmatch(msg.Content)[2]
	r.Sdlr.ResumeWithCommand(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to fishing or hunting event",
	})
}

func (r *Responder) pm(_ discord.Message) {
	res := r.PostmemeOpts[rand.Intn(len(r.PostmemeOpts))]
	r.Sdlr.ResumeWithCommand(&scheduler.Command{
		Value: res,
		Log:   "responding to postmeme",
	})
}

func (r *Responder) event(msg discord.Message) {
	res := exp.event.FindStringSubmatch(msg.Content)[2]
	r.Sdlr.PrioritySchedule(&scheduler.Command{
		Value: clean(res),
		Log:   "responding to global event",
	})
}

func (r *Responder) search(msg discord.Message) {
	choices := exp.search.FindStringSubmatch(msg.Content)[1:]
	for _, choice := range choices {
		for _, allowed := range r.AllowedSearches {
			if choice == allowed {
				r.Sdlr.ResumeWithCommand(&scheduler.Command{
					Value: choice,
					Log:   "responding to search",
				})
				return
			}
		}
	}
	r.Sdlr.ResumeWithCommand(&scheduler.Command{
		Value: r.SearchCancel[rand.Intn(len(r.SearchCancel))],
		Log:   "no allowed search options provided, responding",
	})
}

func (r *Responder) hl(msg discord.Message) {
	if !exp.hl.MatchString(msg.Embeds[0].Description) {
		return
	}
	nstr := strings.Replace(exp.hl.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	n, err := strconv.Atoi(nstr)
	if err != nil {
		r.Logger.Errorf("error while reading highlow hint: %v", err)
		return
	}
	res := "high"
	if n > 50 {
		res = "low"
	}
	r.Sdlr.ResumeWithCommand(&scheduler.Command{
		Value: res,
		Log:   "responding to highlow",
	})
}

func (r *Responder) balCheck(msg discord.Message) {
	if !r.BalanceCheck {
		return
	}
	if !strings.Contains(msg.Embeds[0].Title, r.Client.User.Username) {
		return
	}
	balstr := strings.Replace(exp.bal.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	bal, err := strconv.Atoi(balstr)
	if err != nil {
		r.Logger.Errorf("error while reading balance: %v", err)
		return
	}
	r.Logger.Infof(
		"current wallet balance: %v coins",
		numFmt.Sprintf("%d", bal),
	)
	if r.startingTime.IsZero() {
		r.startingBal = bal
		r.startingTime = time.Now()
		return
	}
	inc := bal - r.startingBal
	per := time.Now().Sub(r.startingTime)
	hourlyInc := int(math.Round(float64(inc) / per.Hours()))
	r.Logger.Infof(
		"average income: %v coins/h",
		numFmt.Sprintf("%d", hourlyInc),
	)
}

func (r *Responder) abLaptop(_ discord.Message) {
	if !r.AutoBuy.Laptop {
		return
	}
	r.Sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy laptop",
		Log:   "no laptop, buying a new one",
	})
}

func (r *Responder) abHuntingRifle(_ discord.Message) {
	if !r.AutoBuy.HuntingRifle {
		return
	}
	r.Sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy rifle",
		Log:   "no hunting rifle, buying a new one",
	})
}

func (r *Responder) abFishingPole(_ discord.Message) {
	if !r.AutoBuy.FishingPole {
		return
	}
	r.Sdlr.PrioritySchedule(&scheduler.Command{
		Value: "pls buy fishing pole",
		Log:   "no fishing pole, buying a new one",
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

func (r *Responder) router() *discord.MessageRouter {
	rtr := &discord.MessageRouter{}

	// Fishing and hunting events.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fh).
		Mentions(r.Client.User.ID).
		Handler(r.fh)

	// Postmeme.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentContains("What type of meme do you want to post").
		Mentions(r.Client.User.ID).
		Handler(r.pm)

	// Global events.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.event).
		Handler(r.event)

	// Search.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.search).
		Mentions(r.Client.User.ID).
		Handler(r.search)

	// Highlow.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		Mentions(r.Client.User.ID).
		Handler(r.hl)

	// Balance report.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		Handler(r.balCheck)

	// Auto buy laptop.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentContains("oi you need to buy a laptop in the shop to post memes").
		Mentions(r.Client.User.ID).
		Handler(r.abLaptop)

	// Auto buy fishing pole.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentContains("You don't have a fishing pole").
		Mentions(r.Client.User.ID).
		Handler(r.abFishingPole)

	// Auto buy hunting rifle.
	rtr.NewRoute().
		Channel(r.ChannelID).
		Author(DMID).
		ContentContains("You don't have a hunting rifle").
		Mentions(r.Client.User.ID).
		Handler(r.abHuntingRifle)

	return rtr
}
