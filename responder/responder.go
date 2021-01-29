// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package responder

import (
	"fmt"
	"time"

	"github.com/dankgrinder/dankgrinder/config"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/scheduler"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Responder struct {
	Sdlr            *scheduler.Scheduler
	Client          *discord.Client
	FatalHandler    func(err error)
	Logger          *logrus.Logger
	ChannelID       string
	PostmemeOpts    []string
	AllowedSearches []string
	BalanceCheck    bool
	AutoBuy         *config.AutoBuy

	ws *discord.WSConn
	startingBal  int
	startingTime time.Time
}

func (r *Responder) Start() error {
	if r.Client == nil {
		return fmt.Errorf("no client")
	}
	if r.Sdlr == nil {
		return fmt.Errorf("no scheduler")
	}
	if r.ChannelID == "" {
		return fmt.Errorf("no channel id")
	}
	if len(r.PostmemeOpts) == 0 {
		return fmt.Errorf("no postmeme options")
	}
	if len(r.AllowedSearches) == 0 {
		return fmt.Errorf("no allowed searches")
	}
	if r.AutoBuy == nil {
		return fmt.Errorf("no auto buy")
	}
	if r.Logger == nil {
		r.Logger = logrus.StandardLogger()
	}
	if r.FatalHandler == nil {
		r.FatalHandler = func(err error) {}
	}

	ws, err := r.Client.NewWSConn(r.router(), r.wsFatalHandler)
	if err != nil {
		return fmt.Errorf("error while starting websocket connection: %v", err)
	}
	r.ws = ws
	r.Logger.Infof("websocket ready")
	return nil
}

func (r *Responder) Close() error {
	if err := r.ws.Close(); err != nil {
		return err
	}
	return nil
}

func (r *Responder) wsFatalHandler(err error) {
	if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == 4004 {
		r.FatalHandler(fmt.Errorf("websocket closed: authentication failed, try using a new token"))
		r.Close()
		return
	}
	r.Logger.Errorf("websocket closed: %v", err)

	r.ws, err = r.Client.NewWSConn(r.router(), r.wsFatalHandler)
	if err != nil {
		r.FatalHandler(fmt.Errorf("error while connecting to websocket: %v", err))
		r.Close()
		return
	}
	r.Logger.Infof("reconnected to websocket")
}
