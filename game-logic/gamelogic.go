package game

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"

	// "github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	LevelChan []chan *Level
	InputChan chan *Input
	Level *Level
}

func NewGame(numWindows int, mapPath string) *Game {
	levelChan := make([]chan *Level, numWindows)
	for i := range levelChan {
		levelChan[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)
	return &Game{levelChan, inputChan, loadMapFromFile(mapPath)}
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
	Search
	CloseWindow
)

type Input struct {
	Input InputState
	X,Y int
	LevelChannel chan *Level
}

type Level struct {
	Map [][]Tile
	Player *Player
	Monsters map[Pos]*Monster
	Debug map[Pos]bool
}

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
	Rune     rune
	Name     string
}

type Character struct {
	Entity
	Hp       int
	Strength int
	Speed    float64
	Ap float64
}

type Player struct{
	Character
}

type AttackInterface interface {
	GetAP() float64
	SetAP(float64)
	GetHP() int
	SetHP(int)
	GetAttackPower() int
}

func (c *Character) GetAP() float64 {
	return c.Ap
}

func (c *Character) SetAP(a float64) {
	c.Ap = a
}

func (c *Character) GetHP() int {
	return c.Hp
}

func (c *Character) SetHP(h int) {
	c.Hp = h
}

func (c *Character) GetAttackPower() int {
	return c.Strength
}

func Attack(a1, a2 AttackInterface) {
	a1.SetAP(a1.GetAP()-1)
	a2.SetHP(a2.GetHP()-a1.GetAttackPower())
	if a2.GetHP() > 0 {
		a2.SetAP(a2.GetAP()-1)
		a1.SetHP(a1.GetHP()-a2.GetAttackPower())
	}
}

// func MonsterAttackPlayer(m *Monster, p *Player) {
// 	m.Ap--
// 	p.Hp -= m.Strength
// 	if p.Hp > 0 {
// 		p.Ap--
// 		m.Hp -= p.Strength
// 	}
// }

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

	level := &Level{Map: make([][]Tile, len(temp)), Monsters: make(map[Pos]*Monster)}

	level.Player = &Player{}
	level.Player.Rune = '@'
	level.Player.Name = "Player"
	level.Player.Hp = 20
	level.Player.Strength = 10
	level.Player.Speed = 1
	level.Player.Ap = 1

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
			} else if col == '@' {
				t = Pending
				level.Player.X = x
				level.Player.Y = y
			} else if col == 'R' {
				t = Pending
				level.Monsters[Pos{x,y}] = NewRat(Pos{x,y})
			} else if col == 'S' {
				t = Pending
				level.Monsters[Pos{x,y}] = NewSpider(Pos{x,y})
			} else {
				panic("the character that you put in map is invalid")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, col := range row {
			if col == Pending {
				level.Map[y][x] = level.bfsFloor(Pos{x,y})
			}
		}
	}

	return level
}

func getNeighbour(level *Level, pos Pos) []Pos {
	res := make([]Pos, 0, 8)
	up := Pos{pos.X,pos.Y-1}
	left := Pos{pos.X-1,pos.Y}
	right := Pos{pos.X+1,pos.Y}
	down := Pos{pos.X,pos.Y+1}
	if canWalk(level,pos.X,pos.Y-1) {
		res = append(res, up)
	}
	if canWalk(level, pos.X-1, pos.Y) {
		res = append(res, left)
	}
	if canWalk(level, pos.X+1, pos.Y) {
		res = append(res, right)
	}
	if canWalk(level, pos.X, pos.Y+1) {
		res = append(res, down)
	}
	return res
}

func (level *Level) bfsFloor(pos Pos) Tile {
	queue := make([]Pos, 0, 8)
	visited := make(map[Pos]bool)
	queue = append(queue, pos)
	visited[pos] = true
	level.Debug = visited

	for len(queue) > 0 {
		curr := queue[0]
		currTile := level.Map[curr.Y][curr.X]
		if currTile == DirtFloor {
			return DirtFloor
		}
		queue = queue[1:]
		for _, adj := range getNeighbour(level, curr) {
			if !visited[adj] && canWalk(level, adj.X, adj.Y) {
				queue = append(queue, adj)
				visited[adj] = true
				time.Sleep(100*time.Millisecond)
			}
		}
	}
	return DirtFloor
}

func (level *Level) aStar(start, goal Pos) []Pos {
	pq := make(priorityQueue, 0, 8)
	pq = pq.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0
	level.Debug = make(map[Pos]bool)

	var curr Pos

	for len(pq) > 0 {
		pq, curr = pq.pop()

		if curr == goal {
			path := make([]Pos, 0, 8)
			p := curr
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)
			for i, j := 0, len(path)-1; i < j; i,j = i+1, j-1 {
				path[i],path[j] = path[j], path[i]
			}
			for _, pos := range path {
				level.Debug[pos] = true
				// time.Sleep(100*time.Millisecond)
			}
			return path
		}

		for _, next := range getNeighbour(level, curr) {
			newCost := costSoFar[curr]+1
			_, exist := costSoFar[next]
			if !exist || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X-next.X)))
				yDist := int(math.Abs(float64(goal.Y-next.Y)))
				priority := newCost+xDist+yDist
				pq = pq.push(next, priority)
				cameFrom[next] = curr
				// level.Debug[next] = true
				// ui.Draw(level)
				// time.Sleep(100*time.Millisecond)
			}
		}
	}
	return nil
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

func (player *Player) move(to Pos, level *Level) {
	monster, exist := level.Monsters[to]
	if !exist {
		player.Pos = to
	} else {
		Attack(level.Player, monster)
		fmt.Println("player hp :", level.Player.Hp)
		fmt.Println("monster hp :", monster.Hp)
		if monster.Hp <= 0 {
			delete(level.Monsters, monster.Pos)
		}
		if level.Player.Hp <= 0 {
			panic("You died!")
		}
	}
} 

func (game *Game) handleInput(input *Input) {
	level := game.Level
	switch input.Input {
	case Up :
		to := Pos{level.Player.X,level.Player.Y-1}
		if isClosedDoor(level,level.Player.X,level.Player.Y-1) {
			level.Map[level.Player.Y-1][level.Player.X] = OpenedDoor
		} else if canWalk(level,level.Player.X,level.Player.Y-1) {
			level.Player.move(to, level)
		}
	case Left :
		to := Pos{level.Player.X-1,level.Player.Y}
		if isClosedDoor(level, level.Player.X-1, level.Player.Y) {
			level.Map[level.Player.Y][level.Player.X-1] = OpenedDoor
		} else if canWalk(level, level.Player.X-1, level.Player.Y) {
			level.Player.move(to, level)
		}
	case Right :
		to := Pos{level.Player.X+1,level.Player.Y}
		if isClosedDoor(level, level.Player.X+1, level.Player.Y) {
			level.Map[level.Player.Y][level.Player.X+1] = OpenedDoor
		} else if canWalk(level, level.Player.X+1, level.Player.Y) {
			level.Player.move(to, level)
		}
	case Down :
		to := Pos{level.Player.X,level.Player.Y+1}
		if isClosedDoor(level, level.Player.X, level.Player.Y+1) {
			level.Map[level.Player.Y+1][level.Player.X] = OpenedDoor
		} else if canWalk(level, level.Player.X, level.Player.Y+1) {
			level.Player.move(to, level)
		}
	case Search:
		// bfs(ui, level, level.Player.Pos)
		level.aStar(level.Player.Pos, Pos{7,4})
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChan {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChan = append(game.LevelChan[:chanIndex], game.LevelChan[chanIndex+1:]...)
	}
}

func (game *Game) Run() {
	for _, lchan := range game.LevelChan {
		lchan <- game.Level
	}

	for input := range game.InputChan {
		if input.Input == Quit {
			return 
		}

		game.handleInput(input)

		for _, monster := range game.Level.Monsters {
			monster.Update(game.Level)
		}

		if len(game.LevelChan) == 0 {
			return
		}

		for _, lchan := range game.LevelChan {
			lchan <- game.Level
		}
	}
}