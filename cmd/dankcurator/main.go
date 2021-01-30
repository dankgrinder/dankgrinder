package main

import (
	"fmt"
	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/mxschmitt/playwright-go"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func boolptr(b bool) *bool {
	return &b
}
func strptr(s string) *string {
	return &s
}
func open(token string) {
	pw, err := playwright.Run()
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: boolptr(false),
	})
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		BypassCSP: boolptr(true),
	})
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	if err = page.AddInitScript(playwright.BrowserContextAddInitScriptOptions{
		Script: strptr(fmt.Sprintf(`
			localStorage.setItem('token', '"%v"');
		`, token)),
	}); err != nil {
		logrus.Fatalf("%v")
	}
	if _, err = page.Goto("https://discord.com/channels/@me"); err != nil {
		logrus.Fatalf("%v")
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	ex, err := os.Executable()
	if err != nil {
		logrus.Fatalf("could not find executable path: %v", err)
	}
	ex = filepath.ToSlash(ex)
	cfg, err := config.Load(path.Dir(ex))
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}
	if err = cfg.Validate(); err != nil {
		logrus.Fatalf("invalid config: %v", err)
	}

	var clients []*discord.Client
	for _, opts := range cfg.InstancesOpts {
		client, err := discord.NewClient(opts.Token)
		if err != nil {
			logrus.Errorf("error while creating client: %v", err)
			continue
		}
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		os.Exit(1)
	}

	for i, client := range clients {
		fmt.Printf("%v#%v (%v)\n", client.User.Username, client.User.Discriminator, i)
	}
	fmt.Printf("choose an instance number to log into: ")
	var s string
	if _, err = fmt.Scanln(&s); err != nil {
		logrus.Fatalf("error while reading stdin: %v", err)
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n >= len(clients) {
		logrus.Fatalf("invalid input")
	}

	open(clients[n].Token)
	<-make(chan bool)
}
