package instance

import (
	"math"
	"math/rand"
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) scratch(msg discord.Message) {

	if strings.Contains(msg.Embeds[0].Description, "You can scratch **0** more fields") {
		in.sdlr.Resume()
		return
	}
	emptyLocations := in.returnButtonEmojiIndex("emptyspace", 3, 5, msg)
	var numberOfCoordinates int = len(emptyLocations) / 2
	p := int(math.Floor(float64(numberOfCoordinates / 2)))
	randomEven := rand.Intn(p) * 2
	actionRow := emptyLocations[randomEven]
	Button := emptyLocations[randomEven+1]
	in.sdlr.ResumeWithCommandOrPrioritySchedule(&scheduler.Command{
		Actionrow:   actionRow,
		Button:      Button,
		Message:     msg,
		Log:         "Pressing random scratch",
		AwaitResume: true,
	})

}

func (in *Instance) scratchEnd(msg discord.Message) {
	in.sdlr.Resume()
	in.sdlr.Resume()
	in.sdlr.Resume()
	in.sdlr.Resume()
}
