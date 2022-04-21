package main

import (

	game "github.com/ahmadfarhanstwn/rpg/game-logic"
	ui2d "github.com/ahmadfarhanstwn/rpg/game-ui-2d"
)

// func init() {

// }

type Ui2D struct {}

func main() {
	ui := &ui2d.Ui2D{}
	game.Run(ui)
}