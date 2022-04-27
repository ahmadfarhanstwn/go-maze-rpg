package game

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	return &Monster{Character{Entity{pos, 'R', "Rat"},10,2,15,0,10}}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{Character{Entity{pos, 'S', "Spider"},15,5,2,0,10}}
}

func (m *Monster) Update(level *Level) {
	m.Ap += m.Speed
	playerPos := level.Player.Pos
	path := level.aStar(m.Pos, playerPos)
	if (len(path) == 0) {
		m.pass()
		return
	}
	apInt := int(m.Ap)
	for i := 1; i <= apInt; i++ {
		if i < len(path) {
			m.Move(path[i], level)
			m.Ap--
		}
	}
}

func (m *Monster) pass() {
	m.Ap -= m.Speed
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exist := level.Monsters[to]
	if !exist && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	} else if to == level.Player.Pos {
		level.Attack(&m.Character, &level.Player.Character)
		level.addEvent(m.Name + "attack player")
		if m.Hp <= 0 {
			level.addEvent("Player killed" + m.Name)
			delete(level.Monsters, m.Pos)
		}
		if level.Player.Hp <= 0 {
			panic("You died!")
		}
	}
}