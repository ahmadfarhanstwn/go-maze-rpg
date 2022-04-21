package game

import (
	"bufio"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type GameUi interface {
	Draw(*Level)
	GetInput() *Input
}

type Tile rune

const (
	StoneWall Tile = '#'
	DirtFloor Tile = '.'
	ClosedDoor Tile = '|'
	OpenedDoor Tile = '/'
	Blank 	  Tile = 0
	Pending   Tile = -1
)

type InputState int

const (
	None InputState = iota
	Left
	Right
	Down
	Up
	Quit
)

type Input struct {
	Input InputState
	X,Y int
}

type Level struct {
	Map [][]Tile
	Player Player
}

type Entity struct {
	X, Y int
}

type Player struct{
	Entity
}

func loadMapFromFile(fileName string) *Level {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	temp := make([]string,0)
	longest, index := 0, 0

	for scanner.Scan() {
		temp = append(temp, scanner.Text())
		if longest < len(temp[index]) {
			longest = len(temp[index])
		}
		index++
	}

	level := &Level{Map: make([][]Tile, len(temp))}

	for y := 0; y < len(level.Map); y++ {
		level.Map[y] = make([]Tile, longest)
		for x, col := range temp[y] {
			var t Tile
			if col == ' ' || col == '\t' || col == '\n' || col == '\r' {
				t = Blank
			} else if col == '#' {
				t = StoneWall
			} else if col == '.' {
				t = DirtFloor
			} else if col == '|' {
				t = ClosedDoor
			} else if col == '/' {
				t = OpenedDoor
			} else if col == 'P' {
				t = Pending
				level.Player.X = x
				level.Player.Y = y
			} else {
				panic("the character that you put in map is invalid")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, col := range row {
			if col == Pending {
				SearchLoop:
					for sy := y-1; sy <= y+1; sy++ {
						for sx := x-1; sx <= x+1; sx++ {
							if level.Map[sy][sx] == DirtFloor {
								level.Map[y][x] = DirtFloor
								break SearchLoop
							}
						}
					}
			}
		}
	}

	return level
}

func isClosedDoor(level *Level, x, y int) bool{
	if x < 0 || x >= int(len(level.Map[0])) || y < 0 || y >= int(len(level.Map)) {
		return false
	}
	return level.Map[y][x] == ClosedDoor
}

func canWalk(level *Level, x, y int) bool {
	if x < 0 || x >= int(len(level.Map[0])) || y < 0 || y >= int(len(level.Map)) {
		return false
	} 
	switch level.Map[y][x] {
	case ClosedDoor,StoneWall, Blank:
		return false
	}
	return true
}

func handleInput(level *Level, input *Input) {
	switch input.Input {
	case Up :
		if isClosedDoor(level,level.Player.X,level.Player.Y-1) {
			level.Map[level.Player.Y-1][level.Player.X] = OpenedDoor
		} else if canWalk(level,level.Player.X,level.Player.Y-1) {
			level.Player.Y--
		}
	case Left :
		if isClosedDoor(level, level.Player.X-1, level.Player.Y) {
			level.Map[level.Player.Y][level.Player.X-1] = OpenedDoor
		} else if canWalk(level, level.Player.X-1, level.Player.Y) {
			level.Player.X--
		}
	case Right :
		if isClosedDoor(level, level.Player.X+1, level.Player.Y) {
			level.Map[level.Player.Y][level.Player.X+1] = OpenedDoor
		} else if canWalk(level, level.Player.X+1, level.Player.Y) {
			level.Player.X++
		}
	case Down :
		if isClosedDoor(level, level.Player.X, level.Player.Y+1) {
			level.Map[level.Player.Y+1][level.Player.X] = OpenedDoor
		} else if canWalk(level, level.Player.X, level.Player.Y+1) {
			level.Player.Y++
		}
	}
}

func Run(gameui GameUi) {
	level := loadMapFromFile("game-logic/maps/level1.map")
	var elapsedTime float32
	for {
		frameStart := time.Now()

		gameui.Draw(level)
		input := gameui.GetInput()
		// do something with input
		if input != nil && input.Input == Quit {
			return
		}

		handleInput(level, input)

		elapsedTime = float32(time.Since(frameStart).Seconds()*1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds()*1000)
		}
	}
}