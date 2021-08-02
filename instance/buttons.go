package instance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
)

func (in *Instance) pressButton(i int, msg discord.Message) {
	url := "https://discord.com/api/v9/interactions"

	data := map[string]interface{}{"component_type": msg.Components[0].Buttons[i].Type, "custom_id": msg.Components[0].Buttons[i].CustomID, "hash": msg.Components[0].Buttons[i].Hash}
	values := map[string]interface{}{"application_id": "270904126974590976", "channel_id": in.ChannelID, "type": "3", "data": data, "guild_id": msg.GuildID, "message_flags": 0, "message_id": msg.ID}
	json_data, err := json.Marshal(values)

	if err != nil {
		fmt.Println(err)
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	req.Header.Set("authorization", in.Client.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func (in *Instance) returnButtonLabel(i int, msg discord.Message) []string {
	buttonLabels := make([]string, i)
	for j := 0; j < i; j++ {
		buttonLabels = append(buttonLabels, msg.Components[0].Buttons[j].Label)
	}
	return buttonLabels
}

func buttonPressDelay(button_press *config.ButtonPressDelay) time.Duration {
	d := button_press.Base
	if button_press.Variation > 0 {
		d += rand.Intn(button_press.Variation)
	}
	return time.Duration(d) * time.Millisecond
}

func (in *Instance) returnButtonIndex(label string, indices int, msg discord.Message) int {
	for k := 0; k < indices; k++ {
		if label == msg.Components[0].Buttons[k].Label {
			return k
		}
	}
	return -1
}
