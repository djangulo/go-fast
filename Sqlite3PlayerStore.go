package poker

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

var ErrRecordAlreadyExists = errors.New("already exists")

func NewSqlite3PlayerStore(file string) (*Sqlite3PlayerStore, func()) {
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}
	db.AutoMigrate(&Player{})

	removeDatabase := func() {
		db.Close()
		os.Remove(file)
	}

	return &Sqlite3PlayerStore{db}, removeDatabase
}

// Sqlite3PlayerStore
type Sqlite3PlayerStore struct {
	DB *gorm.DB
}

func (s *Sqlite3PlayerStore) GetPlayerScore(name string) int {
	var player Player
	s.DB.First(&player)
	return player.Wins
}
func (s *Sqlite3PlayerStore) RecordWin(name string) {
	var player Player
	s.DB.FirstOrCreate(&player, Player{Name: name})
	player.Wins = player.Wins + 1
	s.DB.Save(&player)
}

func (s *Sqlite3PlayerStore) GetLeague() League {
	var players League
	s.DB.Select([]string{"name", "wins"}).Order("wins desc, name").Find(&players)
	return players
}

func (s *Sqlite3PlayerStore) CreatePlayer(player Player) error {
	s.DB.NewRecord(player)
	res := s.DB.Create(&player)
	if res.Error != nil {
		return ErrRecordAlreadyExists
	}
	return nil
}
