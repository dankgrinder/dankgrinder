package instance

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
	"math/rand"

	"github.com/dankgrinder/dankgrinder/discord"
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

	for i := 0; i < len(database.Database); i++ {
		if question == html.UnescapeString(database.Database[i].Question) {
			var res = html.UnescapeString(database.Database[i].Answer)
			i := in.returnButtonIndex(res, 4, msg)
			time.Sleep(1 * time.Second)
			in.pressButton(i, msg)
		}else{
			i := rand.Intn(4)
			in.pressButton(i, msg)
		}
	}
}

