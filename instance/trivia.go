package instance

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

type Database struct {
	Database []TriviaDetail `json:"database"`
}

type TriviaDetail struct {
	Question string `json:"question"`
	Answer   string `json:"correct_answer"`
}

func (in *Instance) trivia(msg discord.Message) {

	details := exp.trivia.FindStringSubmatch(msg.Embeds[0].Description)[1:]
	question := details[0]

	ex, _ := os.Executable()
	ex = filepath.ToSlash(ex)
	jsonFile, err := os.Open((path.Join(path.Dir(ex), "trivia.json")))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	bytevalue, _ := ioutil.ReadAll(jsonFile)

	var database Database
	json.Unmarshal(bytevalue, &database)

	for p := 0; p < len(database.Database); p++ {
		if question == html.UnescapeString(database.Database[p].Question) {
			var res = html.UnescapeString(database.Database[p].Answer)
			p := in.returnButtonIndex(res, 4, msg)

			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Actionrow: 1,
				Button:    p,
				Message:   msg,
				Log:       "Responding to trivia",
			})
			return
		}
	}
	i := rand.Intn(4)
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Actionrow: 1,
		Button: i + 1,
		Message: msg,
		Log: "trivia answer not found, responding with random option",
	})
}
