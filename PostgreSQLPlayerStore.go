package poker

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // unneeded namespace
	"log"
)

var ErrRecordAlreadyExists = errors.New("already exists")
var ErrNoUUID = errors.New("player has no ID (nil uuid.UUID)")

func NewPostgreSQLPlayerStore(host, port, user, dbname, pass string) (*PostgreSQLPlayerStore, func()) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		pass,
		host,
		port,
		dbname,
	)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}
	db.AutoMigrate(&Player{})

	removeDatabase := func() {
		db.Close()
	}

	return &PostgreSQLPlayerStore{db}, removeDatabase
}

// PostgreSQLPlayerStore
type PostgreSQLPlayerStore struct {
	DB *gorm.DB
}

func (s *PostgreSQLPlayerStore) GetPlayerScore(name string) int {
	var player Player
	s.DB.Where("name = ?", name).First(&player)
	return player.Wins
}
func (s *PostgreSQLPlayerStore) RecordWin(name string) {
	var player Player
	s.DB.FirstOrCreate(&player, Player{Name: name})
	player.Wins = player.Wins + 1
	s.DB.Save(&player)
}

func (s *PostgreSQLPlayerStore) GetLeague() League {
	var players League
	s.DB.Select([]string{"name", "wins"}).Order("wins desc, name").Find(&players)
	return players
}

func (s *PostgreSQLPlayerStore) CreatePlayer(player Player) error {
	s.DB.NewRecord(player)
	res := s.DB.Create(&player)
	if res.Error != nil {
		return ErrRecordAlreadyExists
	}
	return nil
}

func (s *PostgreSQLPlayerStore) DeletePlayer(name string) error {
	var player Player
	res := s.DB.First(&player)
	if player.ID.String() == "00000000-0000-0000-0000-000000000000" || player.ID == uuid.Nil || res.Error != nil {
		return ErrNoUUID
	}
	s.DB.Delete(player)
	return nil
}
