// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/dankgrinder/dankgrinder/config"
)

// sec returns n as a time.Duration in seconds.
func sec(n int) time.Duration {
	return time.Duration(n) * time.Second
}

func shiftDur(shift config.Shift) time.Duration {
	if shift.Duration.Base <= 0 {
		return time.Duration(math.MaxInt64)
	}
	d := sec(shift.Duration.Base)
	if shift.Duration.Variation > 0 {
		d += sec(rand.Intn(shift.Duration.Variation))
	}
	return d
}
