package main

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dankMemerID = "270904126974590976"

var conn *discord.WSConn

var (
	searchExp   = regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(?P<option1>.+)\x60,\s\x60(?P<option2>.+)\x60,\s\x60(?P<option3>.+)\x60`)
	fhExp       = regexp.MustCompile(`10\sseconds.*\s?([Tt]yping|[Tt]ype)\s\x60(?P<fh>.+)\x60`)
	cleanExp    = regexp.MustCompile(`[\x20-\x7E]`)
	cleanNumExp = regexp.MustCompile(`[0-9]`)
	hlExp       = regexp.MustCompile(`Your hint is \*\*(?P<n>[0-9]+)\*\*`)
	balExp      = regexp.MustCompile(`\*\*Wallet\*\*: \x60?⏣?\s?(?P<bal>[0-9,]+)\x60?`)
	eventExp    = regexp.MustCompile(`^(Attack the boss by typing|Type) \x60(?P<event>.+)\x60`)
)

// Used to calculate average income.
var initialBal int

// Used to calculate average income.
var startingTime time.Time

var numFormat = message.NewPrinter(language.English)

var noSearchesFound = []string{
	"trash options",
	"tf is this",
	"f off",
	"wth is this",
	"why no good opts?",
	"i dont wanna die",
}

func chatHandler(_ string, msg discord.Message) { // TODO: message update handling.
	if msg.ChannelID != cfg.ChannelID || msg.Author.ID != dankMemerID {
		return
	}
	// Handle fishing and hunting events.
	if fhExp.MatchString(msg.Content) && mentions(msg.Content, user.ID) {
		content := clean(fhExp.FindStringSubmatch(msg.Content)[2], cleanExp)
		logrus.WithField("response", content).Infof("respoding to fishing or hunting event")
		sched.priority <- command{content: content}
		return
	}

	// Handle postmeme.
	if strings.Contains(msg.Content, "What type of meme do you want to post") &&
		mentions(msg.Content, user.ID) {

		content := randElem(cfg.Compat.Postmeme)
		logrus.WithField("response", content).Infof("respoding to postmeme")
		sched.priority <- command{content: content}
		return
	}

	// Handle global events.
	if eventExp.MatchString(msg.Content) {
		content := clean(eventExp.FindStringSubmatch(msg.Content)[2], cleanExp)
		logrus.WithField("response", content).Infof("respoding to global event")
		sched.priority <- command{content: content}
		return
	}

	// Handle search.
	if searchExp.MatchString(msg.Content) && mentions(msg.Content, user.ID) {
		choices := searchExp.FindStringSubmatch(msg.Content)[1:]
		content := chooseSearch(choices)
		logrus.WithField("response", content).Infof("respoding to search")
		sched.priority <- command{content: content}
		return
	}

	// Handle highlow.
	if len(msg.Embeds) > 0 && hlExp.MatchString(msg.Embeds[0].Description) && mentions(msg.Content, user.ID) {
		n, err := strconv.Atoi(hlExp.FindStringSubmatch(msg.Embeds[0].Description)[1])
		if err != nil {
			logrus.Errorf("error while reading highlow hint: %v", err)
			return
		}
		if n > 50 {
			logrus.WithField("response", "low").Infof("respoding to highlow")
			sched.priority <- command{content: "low"}
			return
		}
		logrus.WithField("response", "high").Infof("respoding to highlow")
		sched.priority <- command{content: "high"}
		return
	}

	// Handle balance and report average income.
	if len(msg.Embeds) > 0 &&
		strings.Contains(msg.Embeds[0].Title, fmt.Sprintf("%v's balance", user.Username)) &&
		cfg.Features.BalanceCheck {

		currBal, err := strconv.Atoi(clean(balExp.FindStringSubmatch(msg.Embeds[0].Description)[1], cleanNumExp))
		if err != nil {
			logrus.Errorf("error while reading balance: %v", err)
		}

		if initialBal == 0 {
			logrus.Infof(
				"current wallet balance: %v",
				numFormat.Sprintf("%d", currBal),
			)
			logrus.Infof("no average income available")
			startingTime = time.Now()
			initialBal = currBal
			return
		}

		totalIncome := float64(currBal - initialBal)
		hourlyIncome := int(math.Round(totalIncome / time.Now().Sub(startingTime).Hours()))

		logrus.Infof(
			"current wallet balance: %v",
			numFormat.Sprintf("%d", currBal),
		)
		logrus.Infof(
			"average income: %v coins/h",
			numFormat.Sprintf("%d", hourlyIncome),
		)
		return
	}

	// Respond to no laptop.
	if strings.Contains(msg.Content, "oi you need to buy a laptop in the shop to post memes") &&
		mentions(msg.Content, user.ID) &&
		cfg.Features.AutoBuy.Laptop {

		logrus.WithField("command", "pls buy laptop").Infof("no laptop, buying a new one")
		sched.priority <- command{content: "pls buy laptop"}
		return
	}

	// Respond to no fishing pole.
	if strings.Contains(msg.Content, "You don't have a fishing pole") &&
		mentions(msg.Content, user.ID) &&
		cfg.Features.AutoBuy.FishingPole {

		logrus.WithField("command", "pls buy fishingpole").Infof("no fishing pole, buying a new one")
		sched.priority <- command{content: "pls buy fishingpole"}
		return
	}

	// Respond to no hunting rifle.
	if strings.Contains(msg.Content, "You don't have a hunting rifle") &&
		mentions(msg.Content, user.ID) &&
		cfg.Features.AutoBuy.HuntingRifle {

		logrus.WithField("command", "pls buy rifle").Infof("no hunting rifle, buying a new one")
		sched.priority <- command{content: "pls buy rifle"}
		return
	}
}

func errHandler(err error) {
	logrus.Errorf("websocket error: %v", err)
}

func fatalHandler(err *websocket.CloseError) {
	if err.Code == 4004 {
		logrus.Fatalf("websocket closed: authentication failed, try using a new token")
	}
	logrus.Errorf("websocket closed: %v", err)
	logrus.Infof("reconnecting to websocket")
	connWS()
}

// connWS connects to the Discord websocket. Put in a separate function to avoid
// repetition in fatalHandler.
func connWS() {
	var err error
	conn, err = discord.NewWSConn(cfg.Token, discord.WSConnOpts{
		ChatHandler:  chatHandler,
		ErrHandler:   errHandler,
		FatalHandler: fatalHandler,
	})
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Infof("connected to websocket")
}
