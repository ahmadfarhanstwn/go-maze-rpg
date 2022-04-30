package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

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