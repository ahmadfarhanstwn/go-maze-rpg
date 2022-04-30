package ui2d

import (
	"github.com/ahmadfarhanstwn/rpg/game-logic"
	"github.com/veandco/go-sdl2/sdl"
)

type mouseState struct {
	leftButton  bool
	rightButton bool
	pos         game.Pos
}

func getMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.pos = game.Pos{int(mouseX), int(mouseY)}
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return &result
}

func (ui *ui) keyPressedOnce(key uint8) bool {
	return ui.keyboardState[key] == 1 && ui.prevKeyboardState[key] == 0
}

func (ui *ui) keyPressed(key uint8) bool {
	return ui.keyboardState[key] == 0 && ui.prevKeyboardState[key] == 1
}

func (ui *ui) checkInput(input game.Input) {
	if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
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
		if ui.keyPressedOnce(sdl.SCANCODE_T) {
			input.Input = game.TakeAllItems
		}
		if ui.keyPressedOnce(sdl.SCANCODE_I) {
			if ui.state == MainUI {
				ui.state = InventoryUI
			} else {
				ui.state = MainUI
			}
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
}