package game

import (
	"math/rand"
	"time"
)

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	items := getItemDropped(pos)
	return &Monster{Character{Entity{pos, 'R', "Rat"}, 10, 3, 1, 0, 10, items, nil, nil, nil}}
}

func NewSpider(pos Pos) *Monster {
	// dropped item
	items := getItemDropped(pos)
	return &Monster{Character{Entity{pos, 'S', "Spider"}, 15, 5, 1, 0, 10, items, nil, nil, nil}}
}

func NewGhost(pos Pos) *Monster {
	items := getItemDropped(pos)
	return &Monster{Character{Entity{pos, 'G', "Ghost"}, 20, 10, 1, 0, 10, items, nil, nil, nil}}
}

func getItemDropped(pos Pos) []*Items {
	items := make([]*Items, 0)
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(15)
	switch r {
	case 1:
		items = append(items, newArmour(pos))
	case 2:
		items = append(items, newHelmet(pos))
	case 3:
		items = append(items, newPotion(pos))
	case 4:
		items = append(items, NewSword(pos))
	case 5:
		items = append(items, NewSword(pos), newArmour(pos))
	case 6:
		items = append(items, newHelmet(pos), newPotion(pos))
	}
	return items
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