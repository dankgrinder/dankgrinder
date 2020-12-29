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
	routes []MessageRoute
}

type MessageRoute struct {
	cond    func(msg Message, eventType string) bool
	handler func(msg Message)
}

func (rtr *MessageRouter) process(msg Message, eventType string) {
	for _, rt := range rtr.routes {
		if rt.cond(msg, eventType) {
			rt.handler(msg)
		}
	}
}

func (rt *MessageRoute) addCond(cond func(msg Message, eventType string) bool) *MessageRoute {
	rt.cond = func(msg Message, eventType string) bool {
		return rt.cond(msg, eventType) && cond(msg, eventType)
	}
	return rt
}

func (rtr *MessageRouter) NewRoute() *MessageRoute {
	rt := &MessageRoute{cond: func(_ Message, _ string) bool {
		return true
	}}
	return rt
}

func (rt *MessageRoute) EventType(et string) *MessageRoute {
	return rt.addCond(func(_ Message, eventType string) bool {
		return eventType == et
	})
}

func (rt *MessageRoute) Mentions(id string) *MessageRoute {
	return rt.addCond(func(msg Message, _ string) bool {
		return strings.Contains(msg.Content, fmt.Sprintf("<@%v>", id))
	})
}

func (rt *MessageRoute) ContentMatchesExp(exp *regexp.Regexp) *MessageRoute {
	return rt.addCond(func(msg Message, _ string) bool {
		return exp.MatchString(msg.Content)
	})
}

func (rt *MessageRoute) ContentContains(s string) *MessageRoute {
	return rt.addCond(func(msg Message, _ string) bool {
		return strings.Contains(msg.Content, s)
	})
}

func (rt *MessageRoute) Author(id string) *MessageRoute {
	return rt.addCond(func(msg Message, _ string) bool {
		return msg.Author.ID == id
	})
}

func (rt *MessageRoute) Channel(id string) *MessageRoute {
	return rt.addCond(func(msg Message, _ string) bool {
		return msg.ChannelID == id
	})
}

func (rt *MessageRoute) HasEmbeds(b bool) *MessageRoute {
	return rt.addCond(func(msg Message, _ string) bool {
		if b {
			return len(msg.Embeds) > 0
		}
		return len(msg.Embeds) == 0
	})
}

func (rt *MessageRoute) Handler(h func(msg Message)) {
	rt.handler = h
}
