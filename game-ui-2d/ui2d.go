package ui2d

import (
	"bufio"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ahmadfarhanstwn/rpg/game-logic"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type fontSize int

const (
	smallSize fontSize = iota
	mediumSize
	largeSize
)

type sounds struct {
	footStep []*mix.Chunk 
	openDoor []*mix.Chunk
	attackingSound []*mix.Chunk
}

func playRandomSounds(chunks []*mix.Chunk, volume int) {
	r := rand.Intn(len(chunks))
	chunks[r].Volume(volume)
	chunks[r].Play(-1,0)
}

type ui struct {
	sounds
	winWidth int
	winHeight int
	renderer *sdl.Renderer
	imageAtlas *sdl.Texture
	textureIndex map[rune][]sdl.Rect
	keyboardState []uint8
	prevKeyboardState []uint8
	centerX int
	centerY int
	cameraLimit int
	window *sdl.Window
	levelChannel chan *game.Level
	inputChannel chan *game.Input
	r *rand.Rand
	smallFont *ttf.Font
	mediumFont *ttf.Font
	largeFont *ttf.Font
	texSmall map[string]*sdl.Texture
	texMedium map[string]*sdl.Texture
	texLarge map[string]*sdl.Texture
	eventBackground *sdl.Texture
}

func NewUi(levelChannel chan *game.Level, inputChannel chan *game.Input) *ui {
	ui := &ui{}
	ui.r = rand.New(rand.NewSource(1))
	ui.levelChannel = levelChannel
	ui.inputChannel = inputChannel
	ui.winWidth = 800
	ui.winHeight = 600
	ui.cameraLimit = 5
	window, err := sdl.CreateWindow("RPG BABY!!!", 100, 100,
		int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window
	ui.texSmall = make(map[string]*sdl.Texture)
	ui.texMedium = make(map[string]*sdl.Texture)
	ui.texLarge = make(map[string]*sdl.Texture)

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.imageAtlas = ui.imgFileToTexture("game-ui-2d/assets/tiles.png")

	ui.loadTextureIdx()

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}

	ui.centerX = -1
	ui.centerY = -1

	ui.smallFont, err = ttf.OpenFont("game-ui-2d/assets/Kingthings.ttf", 10)
	if err != nil {
		panic(err)
	}

	ui.mediumFont, err = ttf.OpenFont("game-ui-2d/assets/Kingthings.ttf", 18)
	if err != nil {
		panic(err)
	}

	ui.largeFont, err = ttf.OpenFont("game-ui-2d/assets/Kingthings.ttf", 26)
	if err != nil {
		panic(err)
	}

	ui.eventBackground = ui.GetSinglePixelTex(sdl.Color{0,0,0,128})
	ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
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
	
	mus.Play(-1)

	return ui
}

func (ui *ui) stringToFont(s string, size fontSize, color sdl.Color) *sdl.Texture {
	var font *ttf.Font

	switch size {
	case smallSize:
		tex, exists := ui.texSmall[s]
		if exists {
			return tex
		}
		font = ui.smallFont
	case mediumSize:
		tex, exists := ui.texMedium[s]
		if exists {
			return tex
		}
		font = ui.mediumFont
	case largeSize:
		tex, exists := ui.texLarge[s]
		if exists {
			return tex
		}
		font = ui.largeFont
	}

	fontSurface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic(err)
	}

	fontTexture, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic(err)
	}

	if size == smallSize {
		ui.texSmall[s] = fontTexture
	} else if size == mediumSize {
		ui.texMedium[s] = fontTexture
	} else {
		ui.texLarge[s] = fontTexture
	}

	return fontTexture
}

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	err = mix.Init(mix.INIT_OGG)
}

func(ui *ui) loadTextureIdx() {
	ui.textureIndex = make(map[rune][]sdl.Rect)

	infile, err := os.Open("game-ui-2d/assets/atlas_index.txt")
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := rune(line[0])
		xy := line[1:]
		splitXY := strings.Split(xy, ",")

		x, err := strconv.ParseInt(strings.TrimSpace(splitXY[0]),10,64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXY[1]),10,64)
		if err != nil {
			panic(err)
		}
		count, err := strconv.ParseInt(strings.TrimSpace(splitXY[2]),10,64)
		if err != nil {
			panic(err)
		}

		var rects []sdl.Rect
		for i := 0; i < int(count); i++ {
			rect := sdl.Rect{int32(x*32),int32(y*32),32,32}
			rects = append(rects, rect)
			x++
			if x > 63 {
				x = 0
				y++
			}
		}
		
		ui.textureIndex[tileRune] = rects
	}
}

func(ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func (ui *ui) Draw(level *game.Level) {

	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

	if level.Player.X > ui.centerX+ui.cameraLimit {
		diff := level.Player.X - (ui.centerX+ui.cameraLimit)
		ui.centerX += diff
	} else if level.Player.X < ui.centerX-ui.cameraLimit {
		diff := (ui.centerX-ui.cameraLimit) - level.Player.X
		ui.centerX -= diff
	} else if level.Player.Y > ui.centerY+ui.cameraLimit {
		diff := level.Player.Y - (ui.centerY+ui.cameraLimit)
		ui.centerY += diff
	} else if level.Player.Y < ui.centerY-ui.cameraLimit {
		diff := (ui.centerY-ui.cameraLimit) - level.Player.Y
		ui.centerY -= diff
	}

	offsetX := (ui.winWidth/2) - ui.centerX*32
	offsetY := (ui.winHeight/2) - ui.centerY*32

	ui.renderer.Clear()
	ui.r.Seed(1)
	for y, rows := range level.Map {
		for x, cols := range rows {
			if cols.Rune != game.Blank {
				r := ui.r.Intn(len(ui.textureIndex[cols.Rune]))
				srcRect := ui.textureIndex[cols.Rune][r]
				if level.Map[y][x].Visible || level.Map[y][x].Seen {
					destRect := sdl.Rect{int32(x*32)+int32(offsetX),int32(y*32)+int32(offsetY),32,32}

					pos := game.Pos{x,y}
					if level.Debug[pos] {
						ui.imageAtlas.SetColorMod(128,0,0)
					} else if level.Map[y][x].Seen && !level.Map[y][x].Visible {
						ui.imageAtlas.SetColorMod(128,128,128)
					} else {
						ui.imageAtlas.SetColorMod(255,255,255)
					}
					ui.renderer.Copy(ui.imageAtlas, &srcRect, &destRect)

					if level.Map[y][x].OverlayRune != game.Blank {
						srcRect = ui.textureIndex[cols.OverlayRune][0]
						ui.renderer.Copy(ui.imageAtlas, &srcRect, &destRect)
					}
				}
			}
		}
	}

	ui.imageAtlas.SetColorMod(255,255,255)

	for pos, monster := range level.Monsters {
		monsterRect := ui.textureIndex[monster.Rune][0]
		ui.renderer.Copy(ui.imageAtlas, &monsterRect, &sdl.Rect{int32(pos.X*32)+int32(offsetX),int32(pos.Y*32)+int32(offsetY),32,32})
	}

	playerSrc := ui.textureIndex['@'][0]
	ui.renderer.Copy(ui.imageAtlas, &playerSrc, &sdl.Rect{int32(level.Player.X*32)+int32(offsetX),int32(level.Player.Y*32)+int32(offsetY),32,32})
	
	startEventH := int32(float64(ui.winHeight)*0.75)
	backgroundWidth := int32(float64(ui.winWidth)*0.30)
	
	ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{0, startEventH, backgroundWidth, int32(ui.winHeight)-startEventH})

	for i, event := range level.Events {
		if event != "" {
			tex := ui.stringToFont(event, mediumSize, sdl.Color{255,0,0,0})
			_,_,w,h,err := tex.Query()
			if err != nil {
				panic(err)
			}
			ui.renderer.Copy(tex, nil, &sdl.Rect{0,int32(i*12)+startEventH,w,h})
		}
	}

	hp := ui.stringToFont("Player HP : "+strconv.FormatInt(int64(level.Player.Hp),10),mediumSize, sdl.Color{255,255,255,0})
	_,_,w,h,err := hp.Query()
	if err != nil {
		panic(err)
	}
	ui.renderer.Copy(hp, nil, &sdl.Rect{0,0,w,h})

	ui.renderer.Present()
	// sdl.Delay(16)
}

func (ui *ui) keyPressedOnce(key uint8) bool {
	return ui.keyboardState[key] == 1 && ui.prevKeyboardState[key] == 0
}

func (ui *ui) keyPressed(key uint8) bool {
	return ui.keyboardState[key] == 0 && ui.prevKeyboardState[key] == 1
}

func (ui *ui) GetSinglePixelTex(color sdl.Color) *sdl.Texture {
	tex, _ := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	// if err != nil {
	// 	panic(err)
	// }
	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A
	tex.Update(nil, pixels, 4)
	return tex
}

func (ui *ui) Run() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChannel <- &game.Input{Input: game.Quit}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChannel <- &game.Input{Input: game.CloseWindow,LevelChannel: ui.levelChannel}
				}
			}
		}

		select {
		case newLevel, ok := <-ui.levelChannel:
			if ok {
				switch newLevel.LastEvent {
				case game.Move:
					playRandomSounds(ui.sounds.footStep,20)
				case game.OpenDoor:
					playRandomSounds(ui.openDoor, 75)
				case game.Attacking:
					playRandomSounds(ui.attackingSound, 75)
				}
				ui.Draw(newLevel)
			}
		default:
		}

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			var input game.Input
			if ui.keyPressedOnce(sdl.SCANCODE_UP) || ui.keyPressedOnce(sdl.SCANCODE_W) {
				input.Input = game.Up
			}
			if ui.keyPressedOnce(sdl.SCANCODE_LEFT) || ui.keyPressedOnce(sdl.SCANCODE_A) {
				input.Input = game.Left
			}
			if ui.keyPressedOnce(sdl.SCANCODE_RIGHT) || ui.keyPressedOnce(sdl.SCANCODE_D) {
				input.Input = game.Right
			}
			if ui.keyPressedOnce(sdl.SCANCODE_DOWN) || ui.keyPressedOnce(sdl.SCANCODE_S) {
				input.Input = game.Down
			}
			if ui.keyPressedOnce(sdl.SCANCODE_ESCAPE) {
				input.Input = game.Quit
			}
			// if (ui.keyboardState[sdl.SCANCODE_Q] == 1 && ui.prevKeyboardState[sdl.SCANCODE_Q] == 0) {
			// 	input.Input = game.Search
			// }

			for i := range ui.keyboardState {
				ui.prevKeyboardState[i] = ui.keyboardState[i]
			}

			if input.Input != game.None {
				ui.inputChannel <- &input
			}
		}
		sdl.Delay(10)
	}
}