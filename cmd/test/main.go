package main

import (
	"encoding/json"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	l := launcher.New().Headless(false).Set("window-size", "1280,720").MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect().MustIncognito()
	router := browser.HijackRequests()
	router.MustAdd("https://discord.com/api/v8/auth/login", func(hijack *rod.Hijack) {
		hijack.MustLoadResponse()
		if hijack.Response.Payload().ResponseCode == http.StatusOK {
			body := map[string]string{}
			if err := json.Unmarshal([]byte(hijack.Response.Body()), &body); err != nil {
				logrus.Errorf("error while unmarshalling response body: %v", err)
			}
			logrus.Infof("token: %v", body["token"])
		}
	})
	go router.Run()

	browser.MustPage("https://discord.com/login")
	<-make(chan bool)
}
