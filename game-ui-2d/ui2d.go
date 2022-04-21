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
)

const winWidth, winHeight int = 800, 600

var renderer *sdl.Renderer
var imageAtlas *sdl.Texture
var textureIndex map[game.Tile][]sdl.Rect
var keyboardState []uint8
var prevKeyboardState []uint8
var centerX int
var centerY int

const cameraLimit int = 5

func loadTextureIdx() {
	textureIndex = make(map[game.Tile][]sdl.Rect)

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
		
		textureIndex[tileRune] = rects
		fmt.Println(len(textureIndex[tileRune]))
	}
}

func imgFileToTexture(filename string) *sdl.Texture {
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
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
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

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}

	window, err := sdl.CreateWindow("RPG BABY!!!", 100, 100,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	imageAtlas = imgFileToTexture("game-ui-2d/assets/tiles.png")

	loadTextureIdx()

	keyboardState = sdl.GetKeyboardState()
	prevKeyboardState = make([]uint8, len(keyboardState))
	for i := range keyboardState {
		prevKeyboardState[i] = keyboardState[i]
	}

	centerX = -1
	centerY = -1
}

type Ui2D struct {

}

func (ui *Ui2D) Draw(level *game.Level) {
	if centerX == -1 && centerY == -1 {
		centerX = level.Player.X
		centerY = level.Player.Y
	}

	if level.Player.X > centerX+cameraLimit {
		centerX++
	} else if level.Player.X < centerX-cameraLimit {
		centerX--
	} else if level.Player.Y > centerY+cameraLimit {
		centerY++
	} else if level.Player.Y < centerY-cameraLimit {
		centerY--
	}

	offsetX := (winWidth/2) - centerX*32
	offsetY := (winHeight/2) - centerY*32

	renderer.Clear()
	rand.Seed(1)
	for y, rows := range level.Map {
		for x, cols := range rows {
			if cols != game.Blank {
				r := rand.Intn(len(textureIndex[cols]))
				srcRect := textureIndex[cols][r]
				destRect := sdl.Rect{int32(x*32)+int32(offsetX),int32(y*32)+int32(offsetY),32,32}
				renderer.Copy(imageAtlas, &srcRect, &destRect)
			}
		}
	}
	renderer.Copy(imageAtlas, &sdl.Rect{15*32,94*32,32,32}, &sdl.Rect{int32(level.Player.X*32)+int32(offsetX),int32(level.Player.Y*32)+int32(offsetY),32,32})
	renderer.Present()
}

func (ui *Ui2D) GetInput() *game.Input {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return &game.Input{Input: game.Quit}
			}
		}

		var input game.Input
		if (keyboardState[sdl.SCANCODE_UP] == 0 && prevKeyboardState[sdl.SCANCODE_UP] != 0) || 
		(keyboardState[sdl.SCANCODE_W] == 0 && prevKeyboardState[sdl.SCANCODE_W] != 0) {
			input.Input = game.Up
		}
		if (keyboardState[sdl.SCANCODE_LEFT] == 0 && prevKeyboardState[sdl.SCANCODE_LEFT] != 0) || 
		(keyboardState[sdl.SCANCODE_A] == 0 && prevKeyboardState[sdl.SCANCODE_A] != 0) {
			input.Input = game.Left
		}
		if (keyboardState[sdl.SCANCODE_RIGHT] == 0 && prevKeyboardState[sdl.SCANCODE_RIGHT] != 0) || 
		(keyboardState[sdl.SCANCODE_D] == 0 && prevKeyboardState[sdl.SCANCODE_D] != 0) {
			input.Input = game.Right
		}
		if (keyboardState[sdl.SCANCODE_DOWN] == 0 && prevKeyboardState[sdl.SCANCODE_DOWN] != 0) || 
		(keyboardState[sdl.SCANCODE_S] == 0 && prevKeyboardState[sdl.SCANCODE_S] != 0) {
			input.Input = game.Down
		}
		if (keyboardState[sdl.SCANCODE_ESCAPE] == 0 && prevKeyboardState[sdl.SCANCODE_ESCAPE] != 0) {
			input.Input = game.Quit
		}

		for i := range keyboardState {
			prevKeyboardState[i] = keyboardState[i]
		}

		if input.Input != game.None {
			return &input
		}
	}
}