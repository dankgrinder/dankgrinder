// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (c Config) Validate() error {
	if len(c.Clusters) == 0 {
		return fmt.Errorf("clusters: no clusters, at least 1 is required")
	}
	for ck, cluster := range c.Clusters {
		if err := validateInstance(cluster.Master); err != nil {
			return fmt.Errorf("clusters[%v].master: %v", ck, err)
		}
		for i, instance := range cluster.Instances {
			if err := validateInstance(instance); err != nil {
				return fmt.Errorf("clusters[%v].instances[%v]: %v", ck, i, err)
			}
		}
	}
	if err := validateCompat(c.Compat); err != nil {
		return err
	}
	return nil
}

func validateInstance(instance Instance) error {
	if instance.Token == "" {
		return fmt.Errorf("no token")
	}
	if instance.ChannelID == "" {
		return fmt.Errorf("no channel id")
	}
	if !isValidID(instance.ChannelID) {
		return fmt.Errorf("invalid channel id")
	}
	if len(instance.Shifts) == 0 {
		return fmt.Errorf("no shifts")
	}
	if err := validateFeatures(instance.Features); err != nil {
		return err
	}
	if err := validateShifts(instance.Shifts); err != nil {
		return err
	}
	return nil
}

func validateFeatures(features Features) error {
	if features.AutoSell.Enable {
		if features.AutoSell.Interval < 0 {
			return fmt.Errorf("auto-sell interval must be greater than or equal to 0")
		}
		if len(features.AutoSell.Items) == 0 {
			return fmt.Errorf("auto-sell enabled but no items configured")
		}
	}
	if features.AutoGift.Enable {
		if features.AutoGift.Interval < 0 {
			return fmt.Errorf("auto-gift interval must be greater than or equal to 0")
		}
		if len(features.AutoGift.Items) == 0 {
			return fmt.Errorf("auto-gift enabled but no items configured")
		}
	}
	if features.AutoShare.Enable {
		if features.AutoShare.MinimumBalance < 0 {
			return fmt.Errorf("auto-share minimum must be greater than or equal to 0")
		}
		if features.AutoShare.MaximumBalance < 0 {
			return fmt.Errorf("auto-share maximum must be greater than or equal to 0")
		}
		if features.AutoShare.MinimumBalance > features.AutoShare.MaximumBalance {
			return fmt.Errorf("auto-share minumum must be smaller than or equal to maximum")
		}
	}
	if features.AutoTidepod.Enable && features.AutoTidepod.Interval < 0 {
		return fmt.Errorf("auto-tidepod interval must be greater than or equal to 0")
	}
	if features.BalanceCheck.Enable && features.BalanceCheck.Interval <= 0 {
		return fmt.Errorf("balance check interval must be greater than 0")
	}
	if features.AutoBlackjack.Enable {
		if !features.BalanceCheck.Enable {
			return fmt.Errorf("auto-blackjack enabled but balance check disabled")
		}
		if features.AutoBlackjack.Amount < 0 {
			return fmt.Errorf("auto-blackjack amount must be greater than or equal to 0")
		}
		for colKey, row := range features.AutoBlackjack.LogicTable {
			if colKey != "A" {
				n, err := strconv.Atoi(colKey)
				if err != nil || n < 2 || n > 10 {
					return fmt.Errorf("invalid auto-blackjack logic table key: %v", colKey)
				}
			}
			for rowKey := range row {
				rowKey = strings.Replace(rowKey, "soft", "", -1)
				n, err := strconv.Atoi(rowKey)
				if err != nil || n < 4 || n > 20 {
					return fmt.Errorf("invalid auto-blackjack logic table key: %v", rowKey)
				}
			}
		}
	}

	for i, cmd := range features.CustomCommands {
		if cmd.Value == "" {
			return fmt.Errorf("features.custom_commands[%v].value: no value", i)
		}
		if strings.Contains(cmd.Value, "pls shop") {
			return fmt.Errorf("invalid custom command value: %v, this custom command is disallowed, use auto-gift instead", cmd.Value)
		}
		if strings.Contains(cmd.Value, "pls sell") {
			return fmt.Errorf("invalid custom command value: %v, this custom command is disallowed, use auto-sell instead", cmd.Value)
		}
		if cmd.Amount < 0 {
			return fmt.Errorf("features.custom_commands[%v].amount: value must be greater than or equal to 0", i)
		}
	}
	return nil
}

func validateShifts(shifts []Shift) error {
	for _, shift := range shifts {
		if shift.State != ShiftStateActive && shift.State != ShiftStateDormant {
			return fmt.Errorf("invalid shift state: %v", shift.State)
		}
	}
	return nil
}

func validateCompat(compat Compat) error {
	if len(compat.PostmemeOpts) == 0 {
		return fmt.Errorf("no postmeme compatibility options")
	}
	if len(compat.AllowedSearches) == 0 {
		return fmt.Errorf("no allowed searches")
	}
	if len(compat.AllowedScrambles) == 0 {
		return fmt.Errorf("no allowed scrambles")
	}
	if len(compat.AllowedScramblesWork) == 0 {
		return fmt.Errorf("no allowed work scrambles")
	}
	if len(compat.AllowedFTB) == 0 {
		return fmt.Errorf("no allowed dig fill the blanks")
	}
	if len(compat.DigCancel) == 0 {
		return fmt.Errorf("no dig cancel compatibility options")
	}
	if len(compat.SearchCancel) == 0 {
		return fmt.Errorf("no search cancel compatibility options")
	}
	if len(compat.WorkCancel) == 0 {
		return fmt.Errorf("no work cancel compatibility options")
	}
	if len(compat.AllowedHangman) == 0 {
		return fmt.Errorf("no work hangman options")
	}
	if compat.Cooldown.Dig <= 0 {
		return fmt.Errorf("dig cooldown must be greater than 0")
	}
	if compat.Cooldown.Work <= 0 {
		return fmt.Errorf("work cooldown must be greater than 0")
	}
	if compat.Cooldown.Postmeme <= 0 {
		return fmt.Errorf("postmeme cooldown must be greater than 0")
	}
	if compat.Cooldown.Hunt <= 0 {
		return fmt.Errorf("hunt cooldown must be greater than 0")
	}
	if compat.Cooldown.Highlow <= 0 {
		return fmt.Errorf("highlow cooldown must be greater than 0")
	}
	if compat.Cooldown.Fish <= 0 {
		return fmt.Errorf("fish cooldown must be greater than 0")
	}
	if compat.Cooldown.Search <= 0 {
		return fmt.Errorf("search cooldown must be greater than 0")
	}
	if compat.Cooldown.Beg <= 0 {
		return fmt.Errorf("beg cooldown must be greater than 0")
	}
	if compat.Cooldown.Gift <= 0 {
		return fmt.Errorf("gift cooldown must be greater than 0")
	}
	if compat.Cooldown.Blackjack <= 0 {
		return fmt.Errorf("blackjack cooldown must be greater than 0")
	}
	if compat.Cooldown.Sell <= 0 {
		return fmt.Errorf("sell cooldown must be greater than 0")
	}
	if compat.Cooldown.Share <= 0 {
		return fmt.Errorf("share cooldown must be greater than 0")
	}
	if compat.AwaitResponseTimeout < 0 {
		return fmt.Errorf("await response timeout must be greater than 0")
	}
	return nil
}

func isValidID(id string) bool {
	return regexp.MustCompile(`^[0-9]+$`).Match([]byte(id))
}
