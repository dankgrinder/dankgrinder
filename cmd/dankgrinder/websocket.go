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
	logrus.Errorf("websocket error: %v", err)
}

func fatalHandler(err *websocket.CloseError) {
	if err.Code == 4004 {
		logrus.Fatalf("websocket closed: authentication failed, try using a new token")
	}
	logrus.Errorf("websocket closed: %v", err)
	logrus.Infof("reconnecting to websocket")
	connWS()
}

// connWS connects to the Discord websocket. Put in a separate function to avoid
// repetition in fatalHandler.
func connWS() {
	_, err := discord.NewWSConn(cfg.Token, discord.WSConnOpts{
		MessageRouter: router(),
		ErrHandler:    errHandler,
		FatalHandler:  fatalHandler,
	})
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Infof("connected to websocket")
}
