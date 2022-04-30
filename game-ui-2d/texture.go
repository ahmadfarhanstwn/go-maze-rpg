package ui2d

import (
	"bufio"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) loadTextureIdx() {
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

		x, err := strconv.ParseInt(strings.TrimSpace(splitXY[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXY[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		count, err := strconv.ParseInt(strings.TrimSpace(splitXY[2]), 10, 64)
		if err != nil {
			panic(err)
		}

		var rects []sdl.Rect
		for i := 0; i < int(count); i++ {
			rect := sdl.Rect{int32(x * 32), int32(y * 32), 32, 32}
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

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
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