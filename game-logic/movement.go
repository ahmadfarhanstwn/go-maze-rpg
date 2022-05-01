package game

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
	TakeAllItems
	TakeItem
	DropItem
	EquipItem
)

type Input struct {
	Input        InputState
	Item         *Items
	LevelChannel chan *Level
}

func (level *Level) DropItem(itemToDrop *Items, character *Character) {
	level.LastEvent = DropItems
	level.addEvent("You dropped " + itemToDrop.Name)
	pos := character.Pos
	items := character.Items
	for i, item := range items {
		if item == itemToDrop {
			character.Items = append(character.Items[:i], character.Items[i+1:]...)
			level.Items[pos] = append(level.Items[pos], item)
			return
		}
	}
}

func (level *Level) MoveItem(itemToMove *Items, character *Character) {
	level.LastEvent = PickUpItems
	level.addEvent("You picked " + itemToMove.Name)
	pos := character.Pos
	items := level.Items[pos]
	for i, item := range items {
		if item == itemToMove {
			items = append(items[:i], items[i+1:]...)
			level.Items[pos] = items
			character.Items = append(character.Items, item)
			return
		}
	}
	panic("Tried to move an item we were not on top")
}

func Equip(itemToEquip *Items, character *Character) {
	for i, item := range character.Items {
		if item == itemToEquip {
			character.Items = append(character.Items[:i], character.Items[i+1:]...)
			if item.Type == Helmet {
				character.Helmet = itemToEquip
			} else if item.Type == Sword {
				character.Sword = itemToEquip
			} else if item.Type == Armour {
				character.Armour = itemToEquip
			}
			return
		}
	}
}

func isClosedDoor(level *Level, x, y int) bool {
	if x < 0 || x >= int(len(level.Map[0])) || y < 0 || y >= int(len(level.Map)) {
		return false
	}
	return level.Map[y][x].OverlayRune == ClosedDoor
}

func canWalk(level *Level, x, y int) bool {
	if x < 0 || x >= int(len(level.Map[0])) || y < 0 || y >= int(len(level.Map)) {
		return false
	}
	switch level.Map[y][x].Rune {
	case StoneWall, HellWall, Blank:
		return false
	}
	switch level.Map[y][x].OverlayRune {
	case ClosedDoor:
		return false
	}
	_, exist := level.Monsters[Pos{x, y}]
	return !exist
}

func canSee(level *Level, x, y int) bool {
	if x < 0 || x >= int(len(level.Map[0])) || y < 0 || y >= int(len(level.Map)) {
		return false
	}
	switch level.Map[y][x].Rune {
	case StoneWall, HellWall, Blank:
		return false
	}
	switch level.Map[y][x].OverlayRune {
	case ClosedDoor:
		return false
	}
	return true
}

func (game *Game) move(to Pos) {
	level := game.CurrentLevel
	player := level.Player
	levelandPos := level.Portals[to]
	level.LastEvent = Move
	if levelandPos != nil {
		if level.Coins >= 5 {
			level.LastEvent = Portal
			game.CurrentLevel = levelandPos.Level
			game.CurrentLevel.Player.Pos = levelandPos.Pos
			level.lineOfSight()
		} else {
			level.LastEvent = FailedPortal
		}
	} else {
		player.Pos = to
		for y := range level.Map {
			for x := range level.Map[y] {
				level.Map[y][x].Visible = false
			}
		}
		level.lineOfSight()
	}
}

func (game *Game) resolveMovement(pos Pos) {
	level := game.CurrentLevel
	monster, exist := level.Monsters[pos]
	if exist {
		level.Attack(&level.Player.Character, &monster.Character)
		level.LastEvent = Attacking
		if monster.Character.Hp <= 0 {
			monster.Dead(level)
			level.LastEvent = MonsterDeath
		}
		if level.Player.Character.Hp <= 0 {
			panic("You died")
		}
	} else if canWalk(level, pos.X, pos.Y) {
		game.move(pos) // todo
		if level.Map[pos.Y][pos.X].OverlayRune == Coin {
			level.Coins++
			level.Map[pos.Y][pos.X].OverlayRune = Blank
		}
	} else if isClosedDoor(level, pos.X, pos.Y) {
		level.Map[pos.Y][pos.X].OverlayRune = OpenedDoor
		level.lineOfSight()
		level.LastEvent = OpenDoor
	}
}

func (game *Game) handleInput(input *Input) {
	level := game.CurrentLevel
	switch input.Input {
	case Up:
		to := Pos{level.Player.X, level.Player.Y - 1}
		game.resolveMovement(to) //todo
	case Left:
		to := Pos{level.Player.X - 1, level.Player.Y}
		game.resolveMovement(to)
	case Right:
		to := Pos{level.Player.X + 1, level.Player.Y}
		game.resolveMovement(to)
	case Down:
		to := Pos{level.Player.X, level.Player.Y + 1}
		game.resolveMovement(to)
	case TakeAllItems:
		for _, item := range level.Items[level.Player.Pos] {
			level.MoveItem(item, &level.Player.Character)
		}
	case DropItem:
		level.DropItem(input.Item, &level.Player.Character)
	case TakeItem:
		level.MoveItem(input.Item, &level.Player.Character)
	case EquipItem:
		Equip(input.Item, &level.Player.Character)
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