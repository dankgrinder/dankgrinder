package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Authorization struct {
	Token string
}

type SendMessageOpts struct {
	ChannelID  string
	TypingTime time.Duration

	// If a bool is sent on this channel before the http request is sent (i.e.
	// when it is still typing or just after typing) execution will be aborted
	// and a nil error will be returned.
	Abort chan bool
}

func (auth Authorization) SendMessage(content string, opts SendMessageOpts) error {
	if auth.Token == "" {
		return fmt.Errorf("invalid authorization")
	}
	if opts.ChannelID == "" {
		return fmt.Errorf("no channel id provided")
	}
	if content == "" {
		return fmt.Errorf("no content provided")
	}

	if opts.TypingTime != 0 {
		for i := 0; i < int(opts.TypingTime)/int(time.Second*10); i++ {
			if err := auth.typing(opts.ChannelID); err != nil {
				return err
			}
			select {
			case <-opts.Abort:
				return nil
			case <-time.After(time.Second * 10):
			}
		}
		if err := auth.typing(opts.ChannelID); err != nil {
			return err
		}
		select {
		case <-opts.Abort:
			return nil
		case <-time.After(opts.TypingTime % (time.Second * 10)):
		}
	}

	reqURL := fmt.Sprintf("https://discord.com/api/v8/channels/%v/messages", opts.ChannelID)

	body, err := json.Marshal(&map[string]interface{}{
		"content": content,
		"tts":     false,
	})
	if err != nil {
		return fmt.Errorf("error while encoding message content as json: %v", err)
	}

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	req.Header.Add("Authorization", auth.Token)
	req.Header.Add("User-Agent", "Chrome/86.0.4240.75")
	req.Header.Add("Accept-Language", "en-GB")
	req.Header.Add("Content-Type", "application/json")

	select {
	case <-opts.Abort:
		return nil
	default:
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code while sending message: %v", res.StatusCode)
	}

	return nil
}

func (auth Authorization) CurrentUser() (User, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v8/users/@me", nil)
	if err != nil {
		return User{}, fmt.Errorf("error while creating http request: %v", err)
	}
	auth.headers(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("error while sending http request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("unexpected status code while sending message: %v", res.StatusCode)
	}

	var u User
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		return User{}, fmt.Errorf("error while decoding body: %v", err)
	}
	return u, nil
}

// typing causes Discord to show the "user is typing..." message. It last for 10
// seconds or until the user sends a message in that channel.
//
// Consequently, if you want to make the user type for more than 10 seconds, you
// must call this function every 10 seconds.
func (auth Authorization) typing(channelID string) error {
	reqURL := fmt.Sprintf("https://discord.com/api/v8/channels/%v/typing", channelID)
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	auth.headers(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response code while sending typing message: %v", res.StatusCode)
	}
	return nil
}

func (auth Authorization) headers(r *http.Request) *http.Request {
	r.Header.Add("Authorization", auth.Token)
	r.Header.Add("User-Agent", "Chrome/86.0.4240.75")
	r.Header.Add("Accept-Language", "en-GB")
	return r
}
