package main

import (
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dankMemerID = "270904126974590976"

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

// Used to calculate average income.
var (
	startingBal  int
	startingTime time.Time
)

var numFmt = message.NewPrinter(language.English)

func fh(msg discord.Message) {
	res := exp.fh.FindStringSubmatch(msg.Content)[2]
	sdlr.priority <- &command{
		run: clean(res),
		log: "responding to fishing or hunting event",
	}
}

func pm(_ discord.Message) {
	res := cfg.Compat.Postmeme[rand.Intn(len(cfg.Compat.Postmeme))]
	sdlr.priority <- &command{
		run: res,
		log: "responding to postmeme",
	}
}

func event(msg discord.Message) {
	res := exp.event.FindStringSubmatch(msg.Content)[2]
	sdlr.priority <- &command{
		run: clean(res),
		log: "responding to global event",
	}
}

func search(msg discord.Message) {
	choices := exp.search.FindStringSubmatch(msg.Content)[1:]
	for _, choice := range choices {
		for _, allowed := range cfg.Compat.AllowedSearches {
			if choice == allowed {
				sdlr.priority <- &command{
					run: choice,
					log: "responding to search",
				}
				return
			}
		}
	}
	sdlr.priority <- &command{
		run: []string{
			"trash options",
			"tf is this",
			"f off",
			"wth is this",
			"why no good opts?",
			"i dont wanna die",
		}[rand.Intn(6)], // Update this number to be the length of the slice!
		log: "no allowed search options provided, responding",
	}
}

func hl(msg discord.Message) {
	if !exp.hl.MatchString(msg.Embeds[0].Description) {
		return
	}
	nstr := strings.Replace(exp.hl.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	n, err := strconv.Atoi(nstr)
	if err != nil {
		logrus.StandardLogger().Errorf("error while reading highlow hint: %v", err)
		return
	}
	res := "high"
	if n > 50 {
		res = "low"
	}
	sdlr.priority <- &command{
		run: res,
		log: "responding to highlow",
	}
}

func balCheck(msg discord.Message) {
	if !cfg.Features.BalanceCheck {
		return
	}
	if !strings.Contains(msg.Embeds[0].Title, user.Username) {
		return
	}
	balstr := strings.Replace(exp.bal.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	bal, err := strconv.Atoi(balstr)
	if err != nil {
		logrus.StandardLogger().Errorf("error while reading balance: %v", err)
		return
	}
	logrus.StandardLogger().Infof(
		"current wallet balance: %v coins",
		numFmt.Sprintf("%d", bal),
	)
	if startingTime.IsZero() {
		startingBal = bal
		startingTime = time.Now()
		return
	}
	inc := bal - startingBal
	per := time.Now().Sub(startingTime)
	hourlyInc := int(math.Round(float64(inc) / per.Hours()))
	logrus.StandardLogger().Infof(
		"average income: %v coins/h",
		numFmt.Sprintf("%d", hourlyInc),
	)
}

func abLaptop(_ discord.Message) {
	if !cfg.Features.AutoBuy.Laptop {
		return
	}
	sdlr.priority <- &command{
		run: "pls buy laptop",
		log: "no laptop, buying a new one",
	}
}

func abHuntingRifle(_ discord.Message) {
	if !cfg.Features.AutoBuy.HuntingRifle {
		return
	}
	sdlr.priority <- &command{
		run: "pls buy rifle",
		log: "no hunting rifle, buying a new one",
	}
}

func abFishingPole(_ discord.Message) {
	if !cfg.Features.AutoBuy.FishingPole {
		return
	}
	sdlr.priority <- &command{
		run: "pls buy fishing pole",
		log: "no fishing pole, buying a new one",
	}
}

func router() *discord.MessageRouter {
	rtr := &discord.MessageRouter{}

	// Fishing and hunting events.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentMatchesExp(exp.fh).
		Mentions(user.ID).
		Handler(fh)

	// Postmeme.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentContains("What type of meme do you want to post").
		Mentions(user.ID).
		Handler(pm)

	// Global events.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentMatchesExp(exp.event).
		Handler(event)

	// Search.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentMatchesExp(exp.search).
		Mentions(user.ID).
		Handler(search)

	// Highlow.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		HasEmbeds(true).
		Mentions(user.ID).
		Handler(hl)

	// Balance report.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		HasEmbeds(true).
		Handler(balCheck)

	// Auto buy laptop.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentContains("oi you need to buy a laptop in the shop to post memes").
		Mentions(user.ID).
		Handler(abLaptop)

	// Auto buy fishing pole.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentContains("You don't have a fishing pole").
		Mentions(user.ID).
		Handler(abFishingPole)

	// Auto buy hunting rifle.
	rtr.NewRoute().
		Channel(cfg.ChannelID).
		Author(dankMemerID).
		ContentContains("You don't have a hunting rifle").
		Mentions(user.ID).
		Handler(abHuntingRifle)

	return rtr
}
