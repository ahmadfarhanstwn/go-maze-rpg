package game

type itemtype int

const (
	Sword itemtype = iota
	Armour
	Helmet
	Potion
)

type Items struct {
	Type itemtype
	Entity
	Power float32
}

func NewSword(p Pos) *Items {
	return &Items{Sword,Entity{p, 's',"Sword"},1.5}
}

func newHelmet(p Pos) *Items {
	return &Items{Helmet,Entity{p, 'h', "Helmet"},.2}
}

func newArmour(p Pos) *Items {
	return &Items{Armour,Entity{p, 'a', "Armour"},.3}
}

func newPotion(p Pos) *Items {
	return &Items{Potion,Entity{p, 'p', "Potion"}, 50}
}