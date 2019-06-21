package poker

// Player struct to hold base data
type Player struct {
	ID   int    `sql:"coluw" json:"id"`
	Name string `sql:"column:" json:"name"`
	Wins int    `gorm:"not null;default:0;" json:"wins"`
}
