package main

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
)

// avg returns the average of the passed data.
func avg(data []int) float64 {
	var sum int
	for i := 0; i < len(data); i++ {
		sum += data[i]
	}
	return math.Round(float64(sum) / float64(len(data)))
}

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
		for j := 0; j < len(cfg.Search); j++ {
			if options[i] == cfg.Search[j] {
				return options[i]
			}
		}
	}
	return randElem(noSearchesFound)
}
