package poker

// Player struct to hold base data
type Player struct {
	ID   int    `sql:"column:id;" json:"id"`
	Name string `sql:"column:name;" json:"name"`
	Wins int    `sql:"columns:wins;" json:"wins"`
}
