package game

import (
	"bufio"
	"encoding/csv"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Tile struct {
	Rune rune
	OverlayRune rune
	Visible bool
	Seen bool
}

const (
	StoneWall rune = '#'
	DirtFloor = '.'
	ClosedDoor = '|'
	OpenedDoor = '/'
	Blank = 0
	Pending = -1
	Upstair = 'U'
	Downstair = 'D'
)

type Level struct {
	Map       [][]Tile
	Events    []string
	EventPos  int
	Player    *Player
	Monsters  map[Pos]*Monster
	Portals   map[Pos]*LevelPos
	Items     map[Pos][]*Items
	LastEvent GameEvent
	Debug     map[Pos]bool
}

type LevelPos struct {
	*Level
	Pos
}

type Pos struct {
	X, Y int
}

func (game *Game) loadWorldFile() {
	file, err := os.Open("game-logic/maps/world.txt")
	if err != nil {
		panic(err)
	}
	csvReaders := csv.NewReader(file)
	csvReaders.FieldsPerRecord = -1
	csvReaders.TrimLeadingSpace = true
	rows, err := csvReaders.ReadAll()
	if err != nil {
		panic(err)
	}
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			if game.CurrentLevel == nil {
				panic("couldn't find the name of level in the world file")
			}
			continue
		}
		levelWithPortal := game.Levels[row[0]]
		if levelWithPortal == nil {
			panic("couldn't find the name of level in the world file")
		}
		x, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			panic(err)
		}
		pos := Pos{int(x), int(y)}

		levelToTeleportTo := game.Levels[row[3]]
		if levelToTeleportTo == nil {
			panic("couldn't find the name of level in the world file")
		}
		x, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err = strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			panic(err)
		}
		posToTeleport := Pos{int(x), int(y)}

		levelWithPortal.Portals[pos] = &LevelPos{levelToTeleportTo, posToTeleport}
	}
}

func loadLevels() map[string]*Level {
	player := &Player{}
	player.Rune = '@'
	player.Name = "Player"
	player.Hp = 200
	player.Strength = 5
	player.Speed = 1
	player.Ap = 1
	player.SightRange = 10

	levels := make(map[string]*Level, 0)

	filenames, err := filepath.Glob("game-logic/maps/*.map")
	if err != nil {
		panic(err)
	}
	for _, fileName := range filenames {
		extensionIndex := strings.LastIndex(fileName, ".map")
		lastSlashIndex := strings.LastIndex(fileName, "\\")
		levelName := fileName[lastSlashIndex+1 : extensionIndex]

		file, err := os.Open(fileName)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		temp := make([]string, 0)
		longest, index := 0, 0

		for scanner.Scan() {
			temp = append(temp, scanner.Text())
			if longest < len(temp[index]) {
				longest = len(temp[index])
			}
			index++
		}

		level := &Level{Map: make([][]Tile, len(temp)), Monsters: make(map[Pos]*Monster), Events: make([]string, 12)}

		level.Debug = make(map[Pos]bool)
		level.EventPos = 0
		level.Player = player
		level.Portals = make(map[Pos]*LevelPos)
		level.Items = make(map[Pos][]*Items)

		for y := 0; y < len(level.Map); y++ {
			level.Map[y] = make([]Tile, longest)
			for x, col := range temp[y] {
				var t Tile
				t.OverlayRune = Blank
				if col == ' ' || col == '\t' || col == '\n' || col == '\r' {
					t.Rune = Blank
				} else if col == '#' {
					t.Rune = StoneWall
				} else if col == '.' {
					t.Rune = DirtFloor
				} else if col == '|' {
					t.OverlayRune = ClosedDoor
					t.Rune = Pending
				} else if col == '/' {
					t.Rune = Pending
					t.OverlayRune = OpenedDoor
				} else if col == '@' {
					t.Rune = Pending
					level.Player.X = x
					level.Player.Y = y
				} else if col == 'R' {
					t.Rune = Pending
					level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
				} else if col == 'S' {
					t.Rune = Pending
					level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
				} else if col == 's' {
					t.Rune = Pending
					level.Items[Pos{x, y}] = append(level.Items[Pos{x, y}], NewSword(Pos{x, y}))
				} else if col == 'h' {
					t.Rune = Pending
					level.Items[Pos{x, y}] = append(level.Items[Pos{x, y}], newHelmet(Pos{x, y}))
				} else if col == 'a' {
					t.Rune = Pending
					level.Items[Pos{x, y}] = append(level.Items[Pos{x, y}], newArmour(Pos{x, y}))
				} else if col == 'p' {
					t.Rune = Pending
					level.Items[Pos{x, y}] = append(level.Items[Pos{x, y}], newPotion(Pos{x, y}))
				} else if col == 'D' {
					t.Rune = Pending
					t.OverlayRune = Downstair
				} else if col == 'U' {
					t.Rune = Pending
					t.OverlayRune = Upstair
				} else {
					panic("the character that you put in map is invalid")
				}
				level.Map[y][x] = t
			}
		}

		for y, row := range level.Map {
			for x, col := range row {
				if col.Rune == Pending {
					level.Map[y][x].Rune = level.bfsFloor(Pos{x, y})
				}
			}
		}

		level.lineOfSight()
		levels[levelName] = level
	}

	return levels
}

func getNeighbour(level *Level, pos Pos) []Pos {
	res := make([]Pos, 0, 8)
	up := Pos{pos.X, pos.Y - 1}
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	down := Pos{pos.X, pos.Y + 1}
	if canWalk(level, pos.X, pos.Y-1) {
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

func (level *Level) bfsFloor(pos Pos) rune {
	queue := make([]Pos, 0, 8)
	visited := make(map[Pos]bool)
	queue = append(queue, pos)
	visited[pos] = true

	for len(queue) > 0 {
		curr := queue[0]
		currTile := level.Map[curr.Y][curr.X]
		if currTile.Rune == DirtFloor {
			return DirtFloor
		}
		queue = queue[1:]
		for _, adj := range getNeighbour(level, curr) {
			if !visited[adj] && canWalk(level, adj.X, adj.Y) {
				queue = append(queue, adj)
				visited[adj] = true
				time.Sleep(100 * time.Millisecond)
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
	// level.Debug = make(map[Pos]bool)

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
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			// for _, pos := range path {
			// 	level.Debug[pos] = true
			// 	// time.Sleep(100*time.Millisecond)
			// }
			return path
		}

		for _, next := range getNeighbour(level, curr) {
			newCost := costSoFar[curr] + 1
			_, exist := costSoFar[next]
			if !exist || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
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

func (level *Level) lineOfSight() {
	
	pos := level.Player.Pos
	dist := level.Player.SightRange

	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			d := math.Sqrt(float64(xDelta*xDelta + yDelta*yDelta))
			if d <= float64(dist) {
				level.bresenham(pos, Pos{x, y})
			}
		}
	}
}

func (level *Level) bresenham(start, end Pos) {
	isSteep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))
	if isSteep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}
	deltaY := int(math.Abs(float64(end.Y - start.Y)))
	err := 0
	y := start.Y
	yStep := 1
	if start.Y >= end.Y {
		yStep = -1
	}
	if start.X > end.X {
		deltaX := start.X - end.X
		for x := start.X; x >= end.X; x-- {
			var pos Pos
			if isSteep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Map[pos.Y][pos.X].Visible = true
			level.Map[pos.Y][pos.X].Seen = true
			if !canSee(level, pos.X, pos.Y) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += yStep
				err -= deltaX
			}
		}
	} else {
		deltaX := end.X - start.X
		for x := start.X; x < end.X; x++ {
			var pos Pos
			if isSteep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Map[pos.Y][pos.X].Visible = true
			level.Map[pos.Y][pos.X].Seen = true
			if !canSee(level, pos.X, pos.Y) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += yStep
				err -= deltaX
			}
		}
	}
}

func (level *Level) addEvent(s string) {
	if level.EventPos == len(level.Events) {
		level.Events = level.Events[1:]
		level.Events = append(level.Events, s)
	} else {
		level.Events[level.EventPos] = s
		level.EventPos++
	}
}