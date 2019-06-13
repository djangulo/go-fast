package main

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var ErrRecordAlreadyExists = errors.New("already exists")

func NewSqlite3PlayerStore(file string) *Sqlite3PlayerStore {
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Player{})
	return &Sqlite3PlayerStore{db}
}

type Sqlite3PlayerStore struct {
	db *gorm.DB
}

func (s *Sqlite3PlayerStore) GetPlayerScore(name string) int {
	var player Player
	s.db.First(&player)
	return player.Score
}
func (s *Sqlite3PlayerStore) RecordWin(name string) {
	// fmt.Println("called create")
	var player Player
	s.db.First(&player)
	player.Score = player.Score + 1
	s.db.Save(&player)
	// fmt.Printf("%#v", player)
}

func (s *Sqlite3PlayerStore) CreatePlayer(player Player) error {
	s.db.NewRecord(player)
	res := s.db.Create(&player)
	if res.Error != nil {
		return ErrRecordAlreadyExists
	}
	return nil
}
