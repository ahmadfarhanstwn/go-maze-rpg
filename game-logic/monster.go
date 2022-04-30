package game

import "math/rand"

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	items := make([]*Items, 0)
	r := rand.Intn(2)
	if r == 1 {
		items = append(items, NewSword(pos))
	}
	return &Monster{Character{Entity{pos, 'R', "Rat"}, 10, 2, 15, 0, 10, items, nil, nil, nil}}
}

func NewSpider(pos Pos) *Monster {
	items := make([]*Items, 0)
	r := rand.Intn(2)
	if r == 1 {
		items = append(items, newHelmet(pos))
	}
	return &Monster{Character{Entity{pos, 'S', "Spider"}, 15, 5, 2, 0, 10, items, nil, nil, nil}}
}

func (m *Monster) Update(level *Level) {
	m.Ap += m.Speed
	playerPos := level.Player.Pos
	path := level.aStar(m.Pos, playerPos)
	if len(path) == 0 {
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
			m.Dead(level)
		}
		if level.Player.Hp <= 0 {
			panic("You died!")
		}
	}
}

func (m *Monster) Dead(level *Level) {
	level.addEvent("Player killed" + m.Name)
	delete(level.Monsters, m.Pos)
	groundItems := level.Items[m.Pos]
	for _, item := range m.Items {
		item.Pos = m.Pos
		groundItems = append(groundItems, item)
	}
	level.Items[m.Pos] = groundItems
}