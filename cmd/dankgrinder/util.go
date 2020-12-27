package main

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// randElem returns a random element from the passed slice.
func randElem(elems []string) string {
	return elems[rand.Intn(len(elems))]
}

// clean removes any characters which do not match the expression exp.
func clean(s string, exp *regexp.Regexp) string {
	var c string
	for _, char := range s {
		if exp.MatchString(string(char)) {
			c += string(char)
		}
	}
	return c
}

// mentions returns true if the string contains a Discord mention of the passed
// user id.
func mentions(s string, userID string) bool {
	return strings.Contains(s, fmt.Sprintf("<@%v>", userID))
}

// chooseSearch chooses an option from the "pls search" command which is
// present in the cfg.Search slice. If none are found a random string is
// returned.
func chooseSearch(options []string) string {
	for i := 0; i < len(options); i++ {
		for j := 0; j < len(cfg.Compat.Search); j++ {
			if options[i] == cfg.Compat.Search[j] {
				return options[i]
			}
		}
	}
	return randElem(noSearchesFound)
}

// ms returns n as a time.Duration in milliseconds.
func ms(n int) time.Duration {
	return time.Duration(n) * time.Millisecond
}

// sec returns n as a time.Duration in seconds.
func sec(n int) time.Duration {
	return time.Duration(n) * time.Second
}

func typingTime(cmd string) time.Duration {
	msPerKey := int(math.Round((1.0 / float64(cfg.SuspicionAvoidance.Typing.Speed)) * 60000))

	d := ms(cfg.SuspicionAvoidance.Typing.Base)
	d += ms(len(cmd) * msPerKey)
	if cfg.SuspicionAvoidance.Typing.Variance > 0 {
		d += ms(rand.Intn(cfg.SuspicionAvoidance.Typing.Variance))
	}
	return d
}

func delay() time.Duration {
	d := ms(cfg.SuspicionAvoidance.MessageDelay.Base)
	if cfg.SuspicionAvoidance.MessageDelay.Variance > 0 {
		d += ms(rand.Intn(cfg.SuspicionAvoidance.MessageDelay.Variance))
	}
	return d
}
