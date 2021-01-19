// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func errHandler(err error) {
	logrus.StandardLogger().Errorf("websocket error: %v", err)
}

func fatalHandler(cerr *websocket.CloseError) {
	if cerr.Code == 4004 {
		logrus.StandardLogger().Fatalf("websocket closed: authentication failed, try using a new token")
	}
	logrus.StandardLogger().Errorf("websocket closed: %v", cerr)

	_, err := discord.NewWSConn(cfg.Token, discord.WSConnOpts{
		MessageRouter: router(),
		ErrHandler:    errHandler,
		FatalHandler:  fatalHandler,
	})
	if err != nil {
		logrus.StandardLogger().Fatalf("%v", err)
	}
	logrus.StandardLogger().Infof("reconnected to websocket")
}
