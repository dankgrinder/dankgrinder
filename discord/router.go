// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"regexp"
	"strings"
)

type MessageRouter struct {
	routes []*MessageRoute
}

type MessageRoute struct {
	conds   []condFunc
	handler func(msg Message)
}

type condFunc func(msg Message, eventType string) bool

func (rtr *MessageRouter) process(msg Message, eventType string) {
	for _, rt := range rtr.routes {
		if rt.matches(msg, eventType) {
			rt.handler(msg)
		}
	}
}

func (rt *MessageRoute) matches(msg Message, eventType string) bool {
	for _, cond := range rt.conds {
		if !cond(msg, eventType) {
			return false
		}
	}
	return true
}

func (rtr *MessageRouter) NewRoute() *MessageRoute {
	rt := &MessageRoute{}
	rtr.routes = append(rtr.routes, rt)
	return rt
}

func (rt *MessageRoute) EventType(et string) *MessageRoute {
	rt.conds = append(rt.conds, func(_ Message, eventType string) bool {
		return eventType == et
	})
	return rt
}

func (rt *MessageRoute) Mentions(id string) *MessageRoute {
	rt.conds = append(rt.conds, func(msg Message, _ string) bool {
		return strings.Contains(msg.Content, fmt.Sprintf("<@%v>", id))
	})
	return rt
}

func (rt *MessageRoute) ContentMatchesExp(exp *regexp.Regexp) *MessageRoute {
	rt.conds = append(rt.conds, func(msg Message, _ string) bool {
		return exp.MatchString(msg.Content)
	})
	return rt
}

func (rt *MessageRoute) ContentContains(s string) *MessageRoute {
	rt.conds = append(rt.conds, func(msg Message, _ string) bool {
		return strings.Contains(msg.Content, s)
	})
	return rt
}

func (rt *MessageRoute) Author(id string) *MessageRoute {
	rt.conds = append(rt.conds, func(msg Message, _ string) bool {
		return msg.Author.ID == id
	})
	return rt
}

func (rt *MessageRoute) Channel(id string) *MessageRoute {
	rt.conds = append(rt.conds, func(msg Message, _ string) bool {
		return msg.ChannelID == id
	})
	return rt
}

func (rt *MessageRoute) HasEmbeds(b bool) *MessageRoute {
	rt.conds = append(rt.conds, func(msg Message, _ string) bool {
		if b {
			return len(msg.Embeds) > 0
		}
		return len(msg.Embeds) == 0
	})
	return rt
}

func (rt *MessageRoute) Handler(h func(msg Message)) {
	rt.handler = h
}
