package game

type itemtype int

const (
	Sword itemtype = iota
	Armour
	Helmet
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