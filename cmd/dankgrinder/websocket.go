package main

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/api"
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

var conn *api.WSConn

var (
	searchExp   = regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(?P<option1>.+)\x60,\s\x60(?P<option2>.+)\x60,\s\x60(?P<option3>.+)\x60`)
	fhExp       = regexp.MustCompile(`10\sseconds.*\s?([Tt]yping|[Tt]ype)\s\x60(?P<fh>.+)\x60`)
	cleanExp    = regexp.MustCompile(`[\x20-\x7E]`)
	cleanNumExp = regexp.MustCompile(`[0-9]`)
	hlExp       = regexp.MustCompile(`Your hint is \*\*(?P<n>[0-9]+)\*\*`)
	balExp      = regexp.MustCompile(`\*\*Wallet\*\*: \x60⏣?\s?(?P<bal>[0-9,]+)\x60`)
)

var (
	bal     int   // Used to calculate the income from the past cycle.
	incomes []int // Used to report average coin income, if enabled in the config.
)

var numFormat = message.NewPrinter(language.English)

var noSearchesFound = []string{
	"trash options",
	"tf is this",
	"f off",
	"wth is this",
	"why no good opts?",
	"i dont wanna die",
}

func chatHandler(msg api.Message) {
	if msg.ChannelID != cfg.ChannelID || msg.Author.ID != dankMemerID {
		return
	}
	time.Sleep(time.Duration(cfg.ResDelay) * time.Millisecond)

	// Handle fishing and hunting events.
	if fhExp.MatchString(msg.Content) && mentions(msg.Content, cfg.UserID) {
		content := clean(fhExp.FindStringSubmatch(msg.Content)[2], cleanExp)
		logrus.WithField("response", content).Infof("respoding to fishing or hunting event")
		sendMessage(content)
		return
	}

	// Handle global events.
	for i := 0; i < len(cfg.GlobalEvents); i++ {
		if strings.Contains(clean(msg.Content, cleanExp), fmt.Sprintf("`%v`", cfg.GlobalEvents[i])) {
			logrus.WithField("response", cfg.GlobalEvents[i]).Infof("respoding to global event")
			sendMessage(cfg.GlobalEvents[i])
			return
		}
	}

	// Handle postmeme.
	if strings.Contains(msg.Content, "What type of meme do you want to post") &&
		mentions(msg.Content, cfg.UserID) {

		content := randElem(cfg.Postmeme)
		logrus.WithField("response", content).Infof("respoding to postmeme")
		sendMessage(content)
		return
	}

	// Handle search.
	if searchExp.MatchString(msg.Content) && mentions(msg.Content, cfg.UserID) {
		choices := searchExp.FindStringSubmatch(msg.Content)[1:]
		content := chooseSearch(choices)
		logrus.WithField("response", content).Infof("respoding to search")
		sendMessage(content)
		return
	}

	// Handle highlow.
	if len(msg.Embeds) > 0 && hlExp.MatchString(msg.Embeds[0].Description) && mentions(msg.Content, cfg.UserID) {
		n, err := strconv.Atoi(hlExp.FindStringSubmatch(msg.Embeds[0].Description)[1])
		if err != nil {
			logrus.Errorf("error while reading highlow hint: %v", err)
			return
		}
		if n > 50 {
			logrus.WithField("response", "low").Infof("respoding to highlow")
			sendMessage("low")
			return
		}
		logrus.WithField("response", "high").Infof("respoding to highlow")
		sendMessage("high")
		return
	}

	// Handle balance and report average income.
	if len(msg.Embeds) > 0 &&
		strings.Contains(msg.Embeds[0].Title, fmt.Sprintf("%v's balance", cfg.BalanceCheck.Username)) &&
		cfg.BalanceCheck.Enable {

		currBal, err := strconv.Atoi(clean(balExp.FindStringSubmatch(msg.Embeds[0].Description)[1], cleanNumExp))
		if err != nil {
			logrus.Errorf("error while reading balance: %v", err)
		}

		if bal == 0 {
			logrus.Infof(
				"current wallet balance: %v",
				numFormat.Sprintf("%d", currBal),
			)
			logrus.Infof("no average income available")
			bal = currBal
			return
		}

		income := currBal - bal
		incomes = append(incomes, income)

		// Turn average income per cycle into average income per hour.
		avgIncomeHour := math.Round(avg(incomes) / cycleTime.Hours())

		logrus.Infof(
			"current wallet balance: %v",
			numFormat.Sprintf("%d", currBal),
		)
		logrus.Infof(
			"average income based on %v cycles: %v coins/h",
			numFormat.Sprintf("%d", len(incomes)),
			numFormat.Sprintf("%d", int(avgIncomeHour)),
		)
		bal = currBal
		return
	}

	// Respond to no laptop.
	if strings.Contains(msg.Content, "oi you need to buy a laptop in the shop to post memes") &&
		mentions(msg.Content, cfg.UserID) &&
		cfg.AutoBuy.Laptop {

		logrus.WithField("command", "pls buy laptop").Infof("no laptop, buying a new one")
		sendMessage("pls buy laptop")
		return
	}

	// Respond to no fishing pole.
	if strings.Contains(msg.Content, "You don't have a fishing pole") &&
		mentions(msg.Content, cfg.UserID) &&
		cfg.AutoBuy.FishingPole {

		logrus.WithField("command", "pls buy fishingpole").Infof("no fishing pole, buying a new one")
		sendMessage("pls buy fishingpole")
		return
	}

	// Respond to no hunting rifle.
	if strings.Contains(msg.Content, "You don't have a hunting rifle") &&
		mentions(msg.Content, cfg.UserID) &&
		cfg.AutoBuy.HuntingRifle {

		logrus.WithField("command", "pls buy rifle").Infof("no hunting rifle, buying a new one")
		sendMessage("pls buy rifle")
		return
	}
}

func errHandler(err error) {
	logrus.Errorf("websocket error: %v", err)
}

func fatalHandler(err error) {
	logrus.Errorf("websocket closed: %v", err)
	logrus.Infof("reconnecting to websocket")
	connWS()
}

// connWS connects to the Discord websocket. Put in a separate function to avoid
// repetition in fatalHandler.
func connWS() {
	var err error // Declared beforehand so that conn is a global declaration.
	conn, err = api.NewWSConn(cfg.Token, api.WSConnOpts{
		ChatHandler:  chatHandler,
		ErrHandler:   errHandler,
		FatalHandler: fatalHandler,
	})
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Infof("connected to websocket")
}
