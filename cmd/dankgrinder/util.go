// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"math"
	"math/rand"
	"regexp"
	"time"
)

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

// ms returns n as a time.Duration in milliseconds.
func ms(n int) time.Duration {
	return time.Duration(n) * time.Millisecond
}

// sec returns n as a time.Duration in seconds.
func sec(n int) time.Duration {
	return time.Duration(n) * time.Second
}

// typing returns a duration for which to type based on the variables in the
// config.
func typing(cmd string) time.Duration {
	msPerKey := int(math.Round((1.0 / float64(cfg.SuspicionAvoidance.Typing.Speed)) * 60000))
	d := ms(cfg.SuspicionAvoidance.Typing.Base)
	d += ms(len(cmd) * msPerKey)
	if cfg.SuspicionAvoidance.Typing.Variance > 0 {
		d += ms(rand.Intn(cfg.SuspicionAvoidance.Typing.Variance))
	}
	return d
}

// delay returns a duration for which to sleep before commencing typing based on
// the variables in the config.
func delay() time.Duration {
	d := ms(cfg.SuspicionAvoidance.MessageDelay.Base)
	if cfg.SuspicionAvoidance.MessageDelay.Variance > 0 {
		d += ms(rand.Intn(cfg.SuspicionAvoidance.MessageDelay.Variance))
	}
	return d
}
