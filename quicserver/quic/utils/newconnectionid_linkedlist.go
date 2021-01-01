package utils

type NewConnectionIDElement struct {
	new, prev *NewConnectionIDElement
	list      *NewConnectionIDElement
	Value     NewConnectionID
}

type NewConnectionIDList struct {
	root NewConnectionIDElement
	len  int
}
