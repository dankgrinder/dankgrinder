package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SendMessageOpts struct {
	Token     string
	ChannelID string
	Content   string
	Typing    time.Duration
}

func SendMessage(opts SendMessageOpts) error {
	if opts.Token == "" {
		return fmt.Errorf("no token provided")
	}
	if opts.ChannelID == "" {
		return fmt.Errorf("no channel id provided")
	}
	if opts.Content == "" {
		return fmt.Errorf("no content provided")
	}

	if opts.Typing != 0 {
		for i := 0; i < int(opts.Typing)/int(time.Second*10); i++ {
			if err := typing(opts.Token, opts.ChannelID); err != nil {
				return err
			}
			time.Sleep(time.Second * 10)
		}
		if err := typing(opts.Token, opts.ChannelID); err != nil {
			return err
		}
		time.Sleep(opts.Typing % (time.Second * 10))
	}

	reqURL := fmt.Sprintf("https://discord.com/api/v8/channels/%v/messages", opts.ChannelID)

	body, err := json.Marshal(&map[string]interface{}{
		"content": opts.Content,
		"tts":     false,
	})
	if err != nil {
		return fmt.Errorf("error while encoding message content as json: %v", err)
	}

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	req.Header.Add("Authorization", opts.Token)
	req.Header.Add("User-Agent", "Chrome/86.0.4240.75")
	req.Header.Add("Accept-Language", "en-GB")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code while sending message: %v", res.StatusCode)
	}

	return nil
}

func typing(token, channelID string) error {
	reqURL := fmt.Sprintf("https://discord.com/api/v8/channels/%v/typing", channelID)
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("User-Agent", "Chrome/86.0.4240.75")
	req.Header.Add("Accept-Language", "en-GB")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response code while sending typing message: %v", res.StatusCode)
	}
	return nil
}
