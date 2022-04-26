package ui2d

import (
	"bufio"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ahmadfarhanstwn/rpg/game-logic"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type ui struct {
	winWidth int
	winHeight int
	renderer *sdl.Renderer
	imageAtlas *sdl.Texture
	textureIndex map[game.Tile][]sdl.Rect
	keyboardState []uint8
	prevKeyboardState []uint8
	centerX int
	centerY int
	cameraLimit int
	window *sdl.Window
	levelChannel chan *game.Level
	inputChannel chan *game.Input
	r *rand.Rand
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

	font, err := ttf.OpenFont("game-ui-2d/assets/Kingthing.ttf", 32)
	if err != nil {
		panic(err)
	}

	fontSurface, err := font.RenderUTF8Solid("Hello", sdl.Color{255,0,0,0})
	if err != nil {
		panic(err)
	}

	fontTexture, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic(err)
	}

	return ui
}

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func(ui *ui) loadTextureIdx() {
	ui.textureIndex = make(map[game.Tile][]sdl.Rect)

	infile, err := os.Open("game-ui-2d/assets/atlas_index.txt")
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := game.Tile(line[0])
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
		ui.centerX++
	} else if level.Player.X < ui.centerX-ui.cameraLimit {
		ui.centerX--
	} else if level.Player.Y > ui.centerY+ui.cameraLimit {
		ui.centerY++
	} else if level.Player.Y < ui.centerY-ui.cameraLimit {
		ui.centerY--
	}

	offsetX := (ui.winWidth/2) - ui.centerX*32
	offsetY := (ui.winHeight/2) - ui.centerY*32

	ui.renderer.Clear()
	ui.r.Seed(1)
	for y, rows := range level.Map {
		for x, cols := range rows {
			if cols != game.Blank {
				r := ui.r.Intn(len(ui.textureIndex[cols]))
				srcRect := ui.textureIndex[cols][r]
				destRect := sdl.Rect{int32(x*32)+int32(offsetX),int32(y*32)+int32(offsetY),32,32}

				pos := game.Pos{x,y}
				if level.Debug[pos] {
					ui.imageAtlas.SetColorMod(128,0,0)
				} else {
					ui.imageAtlas.SetColorMod(255,255,255)
				}
				ui.renderer.Copy(ui.imageAtlas, &srcRect, &destRect)
			}
		}
	}
	for pos, monster := range level.Monsters {
		monsterRect := ui.textureIndex[game.Tile(monster.Rune)][0]
		ui.renderer.Copy(ui.imageAtlas, &monsterRect, &sdl.Rect{int32(pos.X*32)+int32(offsetX),int32(pos.Y*32)+int32(offsetY),32,32})
	}

	playerSrc := ui.textureIndex['@'][0]
	ui.renderer.Copy(ui.imageAtlas, &playerSrc, &sdl.Rect{int32(level.Player.X*32)+int32(offsetX),int32(level.Player.Y*32)+int32(offsetY),32,32})
	ui.renderer.Present()
	// sdl.Delay(16)
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
				ui.Draw(newLevel)
			}
		default:
		}

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			var input game.Input
			if (ui.keyboardState[sdl.SCANCODE_UP] == 1 && ui.prevKeyboardState[sdl.SCANCODE_UP] == 0) || 
			(ui.keyboardState[sdl.SCANCODE_W] == 1 && ui.prevKeyboardState[sdl.SCANCODE_W] == 0) {
				input.Input = game.Up
			}
			if (ui.keyboardState[sdl.SCANCODE_LEFT] == 1 && ui.prevKeyboardState[sdl.SCANCODE_LEFT] == 0) || 
			(ui.keyboardState[sdl.SCANCODE_A] == 1 && ui.prevKeyboardState[sdl.SCANCODE_A] == 0) {
				input.Input = game.Left
			}
			if (ui.keyboardState[sdl.SCANCODE_RIGHT] == 1 && ui.prevKeyboardState[sdl.SCANCODE_RIGHT] == 0) || 
			(ui.keyboardState[sdl.SCANCODE_D] == 1 && ui.prevKeyboardState[sdl.SCANCODE_D] == 0) {
				input.Input = game.Right
			}
			if (ui.keyboardState[sdl.SCANCODE_DOWN] == 1 && ui.prevKeyboardState[sdl.SCANCODE_DOWN] == 0) || 
			(ui.keyboardState[sdl.SCANCODE_S] == 1 && ui.prevKeyboardState[sdl.SCANCODE_S] == 0) {
				input.Input = game.Down
			}
			if (ui.keyboardState[sdl.SCANCODE_ESCAPE] == 1 && ui.prevKeyboardState[sdl.SCANCODE_ESCAPE] == 0) {
				input.Input = game.Quit
			}
			if (ui.keyboardState[sdl.SCANCODE_Q] == 1 && ui.prevKeyboardState[sdl.SCANCODE_Q] == 0) {
				input.Input = game.Search
			}

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