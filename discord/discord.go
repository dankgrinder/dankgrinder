// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUnauthorized    = fmt.Errorf("invalid authorization, try using a new token")
	ErrForbidden       = fmt.Errorf("forbidden, you may not have permission to send in the channel (i.e. you aren't in the server or don't have send message permissions in the channel), your account might need verification, or your ip address may have been blocked")
	ErrTooManyRequests = fmt.Errorf("you are being rate limited, try waiting some time and trying again")
	ErrNotFound        = fmt.Errorf("not found, make sure your channel id is valid")
	ErrIntervalServer  = fmt.Errorf("remote server interval server error")
)

type Client struct {
	Token string
	User  User
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("no token")
	}
	c := &Client{Token: token}
	u, err := c.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("could not get user information: %v", err)
	}
	c.User = u
	return c, nil
}

func (client Client) SendMessage(content, channelID string, typing time.Duration) error {
	if client.Token == "" {
		return fmt.Errorf("no token")
	}
	if channelID == "" {
		return fmt.Errorf("no channel id")
	}
	if content == "" {
		return fmt.Errorf("no content")
	}

	if typing != 0 {
		iterations := int(int64(typing)/int64(time.Second*10)) + 1
		for i := 0; i < iterations; i++ {
			if err := client.typing(channelID); err != nil {
				return err
			}
			s := time.Second * 10
			if i == iterations-1 { // If this is the last iteration.
				s = typing % (time.Second * 10)
			}
			time.Sleep(s)
		}
	}

	reqURL := fmt.Sprintf("https://discord.com/api/v9/channels/%v/messages", channelID)

	body, err := json.Marshal(&map[string]interface{}{
		"content": content,
		"tts":     false,
		"nonce":   client.snowflake(),
	})
	if err != nil {
		return fmt.Errorf("error while encoding message content as json: %v", err)
	}

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	req.Header.Add("Authorization", client.Token)
	req.Header.Add("User-Agent", "Chrome/94.0.4606.31")
	req.Header.Add("Accept-Language", "en-GB")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusForbidden:
			return ErrForbidden
		case http.StatusNotFound:
			return ErrNotFound
		case http.StatusTooManyRequests:
			return ErrTooManyRequests
		case http.StatusInternalServerError:
			return ErrIntervalServer
		default:
			return fmt.Errorf("unexpected status code while sending message: %v", res.StatusCode)
		}
	}
	return nil
}
func (client Client) PressButton(i int, k int, msg Message) error {
	if client.Token == "" {
		return fmt.Errorf("no token")
	}

	i--
	k--

	x := rand.Intn(500)

	time.Sleep(time.Duration(x) * time.Millisecond)

	url := "https://discord.com/api/v9/interactions"

	data := map[string]interface{}{"component_type": msg.Components[i].Buttons[k].Type, "custom_id": msg.Components[i].Buttons[k].CustomID, "hash": msg.Components[i].Buttons[k].Hash}
	values := map[string]interface{}{"application_id": 270904126974590976, "channel_id": msg.ChannelID, "type": "3", "data": data, "guild_id": msg.GuildID, "message_flags": 0, "message_id": msg.ID, "nonce": client.snowflake()}
	json_data, err := json.Marshal(values)

	if err != nil {
		return fmt.Errorf("error while encoding button click as json: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	req.Header.Set("authorization", client.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Chrome/94.0.4606.31 ")
	req.Header.Set("Accept-Language", "en-GB")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}
	if resp.StatusCode != 204 {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusForbidden:
			return ErrForbidden
		case http.StatusNotFound:
			return ErrNotFound
		case http.StatusTooManyRequests:
			return ErrTooManyRequests
		case http.StatusInternalServerError:
			return ErrIntervalServer
		case 400:
			return nil
		default:
			return fmt.Errorf("unexpected status code while clicking button: %v", resp.StatusCode)
		}
	}
	return nil
}

// CurrentUser sends a http request to Discord and returns a User struct based
// on the response. This method should only be used if the user information was
// changed between when you created the client and now. Otherwise, this is also
// available in the User field of the Client struct.
func (client Client) CurrentUser() (User, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v8/users/@me", nil)
	if err != nil {
		return User{}, fmt.Errorf("error while creating http request: %v", err)
	}
	client.headers(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("error while sending http request: %v", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return User{}, ErrUnauthorized
	}
	if res.StatusCode == http.StatusForbidden {
		return User{}, ErrForbidden
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
func (client Client) typing(channelID string) error {
	if channelID == "" {
		return fmt.Errorf("no channel id")
	}
	reqURL := fmt.Sprintf("https://discord.com/api/v8/channels/%v/typing", channelID)
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error while creating http request: %v", err)
	}
	client.headers(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending http request: %v", err)
	}
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusForbidden:
			return ErrForbidden
		case http.StatusNotFound:
			return ErrNotFound
		case http.StatusTooManyRequests:
			return ErrTooManyRequests
		default:
			return fmt.Errorf("unexpected status code while sending message: %v", res.StatusCode)
		}
	}
	return nil
}

func (client Client) headers(r *http.Request) *http.Request {
	r.Header.Add("Authorization", client.Token)
	r.Header.Add("User-Agent", "Chrome/86.0.4240.75")
	r.Header.Add("Accept-Language", "en-GB")
	return r
}

func (client Client) snowflake() int64 {
	snowflake := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(snowflake, 2, 64)
	return nonce
}
