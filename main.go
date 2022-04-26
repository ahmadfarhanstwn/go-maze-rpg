package main

import (
	"runtime"

	game "github.com/ahmadfarhanstwn/rpg/game-logic"
	ui2d "github.com/ahmadfarhanstwn/rpg/game-ui-2d"
)

func main() {
	game := game.NewGame(1, "game-logic/maps/level1.map")
	for i := 0; i < 1; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui := ui2d.NewUi(game.LevelChan[i], game.InputChan)
			ui.Run()
		}(i)
	}
	game.Run()
}