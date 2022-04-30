//TODO :
// ADD MORE MONSTERS
// ADD MORE MAPS
// ADD MORE ITEMS
// ADD MORE TILES
// IMPROVE UI DESIGN FOR INVENTORY
// ADD MORE GAMEPLAY

package main

import (
	"runtime"

	game "github.com/ahmadfarhanstwn/rpg/game-logic"
	ui2d "github.com/ahmadfarhanstwn/rpg/game-ui-2d"
)

func main() {
	game := game.NewGame(1)
	for i := 0; i < 1; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui := ui2d.NewUi(game.LevelChan[i], game.InputChan)
			ui.Run()
		}(i)
	}
	game.Run()
}