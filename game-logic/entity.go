package game

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
	SightRange int
	Items []*Items
	Sword, Armour, Helmet *Items
}

type Player struct{
	Character
}

func (level *Level) Attack(c1, c2 *Character) {
	c1.Ap--
	c1AP := c1.Strength

	if c1.Sword != nil {
		c1AP *= int(c1.Sword.Power)
	}
	if c2.Helmet != nil {
		c1AP = int(float32(c1AP)*(1.0-c2.Helmet.Power))
	}
	if c2.Armour != nil {
		c1AP = int(float32(c1AP)*(1.0-c2.Armour.Power))
	}
	c2.Hp -= c1AP

	if c2.Hp > 0 {
		level.addEvent(c1.Name + " attacked " + c2.Name)
	} else {
		level.addEvent(c1.Name + " killed " + c2.Name)
	}
}