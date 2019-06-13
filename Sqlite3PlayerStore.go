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
	return player.Wins
}
func (s *Sqlite3PlayerStore) RecordWin(name string) {
	var player Player
	s.db.FirstOrCreate(&player, Player{Name: name})
	player.Wins = player.Wins + 1
	s.db.Save(&player)
}

func (s *Sqlite3PlayerStore) GetLeague() []Player {
	var players []Player
	s.db.Select([]string{"name", "wins"}).Find(&players)
	return players
}

func (s *Sqlite3PlayerStore) CreatePlayer(player Player) error {
	s.db.NewRecord(player)
	res := s.db.Create(&player)
	if res.Error != nil {
		return ErrRecordAlreadyExists
	}
	return nil
}
