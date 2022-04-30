package game

type GameEvent int

const (
	Move GameEvent = iota
	OpenDoor
	Portal
	Attacking
	PickUpItems
	DropItems
	EquipItems
	MonsterDeath
	DrinkPotion
)

type Game struct {
	LevelChan []chan *Level
	InputChan chan *Input
	Levels map[string]*Level
	CurrentLevel *Level
}

func NewGame(numWindows int) *Game {
	levelChan := make([]chan *Level, numWindows)
	for i := range levelChan {
		levelChan[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)
	levels := loadLevels()
	
	game := &Game{levelChan, inputChan, levels, nil}
	game.loadWorldFile()
	game.CurrentLevel.lineOfSight()
	return game
}

func (game *Game) Run() {
	for _, lchan := range game.LevelChan {
		lchan <- game.CurrentLevel
	}

	for input := range game.InputChan {
		if input.Input == Quit {
			return 
		}

		p := game.CurrentLevel.Player.Pos
		game.CurrentLevel.bresenham(p, Pos{p.X+7,p.Y-3})
		// for _, pos := range bres {
		// 	game.Level.Debug[pos] = true
		// }

		game.handleInput(input)

		for _, monster := range game.CurrentLevel.Monsters {
			monster.Update(game.CurrentLevel)
		}

		if len(game.LevelChan) == 0 {
			return
		}

		for _, lchan := range game.LevelChan {
			lchan <- game.CurrentLevel
		}
	}
}