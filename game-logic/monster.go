package game

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	monster := &Monster{}
	monster.Pos = pos
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hp = 10
	monster.Strength = 2
	monster.Speed = 15
	monster.Ap = 0
	return monster
}

func NewSpider(pos Pos) *Monster {
	monster := &Monster{}
	monster.Pos = pos
	monster.Rune = 'S'
	monster.Name = "Spider"
	monster.Hp = 15
	monster.Strength = 5
	monster.Speed = 2
	monster.Ap = 0
	return monster
}

func (m *Monster) Update(level *Level) {
	m.Ap += m.Speed
	playerPos := level.Player.Pos
	path := level.aStar(m.Pos, playerPos)
	apInt := int(m.Ap)
	for i := 1; i <= apInt; i++ {
		if i < len(path) {
			m.Move(path[i], level)
			m.Ap--
		}
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exist := level.Monsters[to]
	if !exist && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	} else {
		Attack(m, level.Player)
		if m.Hp <= 0 {
			delete(level.Monsters, m.Pos)
		}
		if level.Player.Hp <= 0 {
			panic("You died!")
		}
	}
}