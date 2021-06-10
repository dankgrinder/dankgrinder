package instance

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
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
	choices := map[string]string{details[2]: details[1], details[4]: details[3], details[6]: details[5]}

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

	for i := 0; i < len(database.Database); i++ {
		if question == html.UnescapeString(database.Database[i].Question) { // changes
			var res = choices[database.Database[i].Answer]
			in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
				Value: res,
				Log:   "responding to Work Color",
			})
		}

	}

}
