package ui2d

import (
	"math/rand"
	"strconv"

	"github.com/veandco/go-sdl2/mix"
)

type sounds struct {
	footStep        []*mix.Chunk
	openDoor        []*mix.Chunk
	attackingSound  []*mix.Chunk
	pickUpItems     []*mix.Chunk
	enteringPortals []*mix.Chunk
}

func (ui *ui) initSound() {
	err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}

	mus, err := mix.LoadMUS("game-ui-2d/assets/cave_theme.ogg")
	if err != nil {
		panic(err)
	}

	stepBase := "game-ui-2d/assets/stepdirt_"
	for i := 1; i <= 8; i++ {
		soundFile := stepBase + strconv.Itoa(i) + ".wav"
		stepSound, err := mix.LoadWAV(soundFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.footStep = append(ui.sounds.footStep, stepSound)
	}

	openDoorSound, err := mix.LoadWAV("game-ui-2d/assets/open_door.ogg")
	if err != nil {
		panic(err)
	}
	ui.sounds.openDoor = append(ui.sounds.openDoor, openDoorSound)

	stepBase = "game-ui-2d/assets/Mudchute_pig_"
	for i := 1; i <= 3; i++ {
		soundFile := stepBase + strconv.Itoa(i) + ".ogg"
		stepSound, err := mix.LoadWAV(soundFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.attackingSound = append(ui.sounds.attackingSound, stepSound)
	}

	pickUpItem, err := mix.LoadWAV("game-ui-2d/assets/metal-clash.wav")
	if err != nil {
		panic(err)
	}

	ui.sounds.pickUpItems = append(ui.sounds.pickUpItems, pickUpItem)

	portalSound, err := mix.LoadWAV("game-ui-2d/assets/porta.ogg")
	if err != nil {
		panic(err)
	}

	ui.sounds.enteringPortals = append(ui.sounds.enteringPortals, portalSound)
	
	mus.Play(-1)
}

func playRandomSounds(chunks []*mix.Chunk, volume int) {
	r := rand.Intn(len(chunks))
	chunks[r].Volume(volume)
	chunks[r].Play(-1,0)
}