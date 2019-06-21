package poker

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq" // unneeded namespace
	"log"
)

var ErrRecordAlreadyExists = errors.New("already exists")

func NewPostgreSQLPlayerStore(host, port, user, dbname, pass string) (*PostgreSQLPlayerStore, func()) {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user,
		pass,
		host,
		port,
		dbname,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	_, errCreate := db.Exec(`
	CREATE TABLE IF NOT EXISTS players (
		id		serial		PRIMARY KEY,
		name	varchar(80)	NOT NULL UNIQUE,
		wins	int			DEFAULT 0
	);
	`)
	if errCreate != nil {
		log.Fatalf("failed to create table %v", errCreate)
	}

	removeDatabase := func() {
		db.Close()
	}

	return &PostgreSQLPlayerStore{db}, removeDatabase
}

// PostgreSQLPlayerStore
type PostgreSQLPlayerStore struct {
	DB *sql.DB
}

func (s *PostgreSQLPlayerStore) GetPlayerScore(name string) int {
	var wins int
	err := s.DB.QueryRow(`SELECT wins FROM players WHERE name = $1;`, name).Scan(&wins)
	if err != nil {
		log.Printf("error: %v", err)
		return 0
	}
	return wins
}
func (s *PostgreSQLPlayerStore) RecordWin(name string) {
	var userID int
	err := s.DB.QueryRow(`SELECT id FROM players WHERE name = $1;`, name).Scan(&userID)
	if err != nil { // likely does not exist
		log.Printf("error: %v, inserting", err)
		s.DB.Exec(`INSERT INTO players(name, wins) VALUES($1, 1);`, name)
		return
	}
	s.DB.Exec(`UPDATE players SET wins = wins + 1 WHERE name = $1`, name)
}

func (s *PostgreSQLPlayerStore) GetLeague() League {
	results, err := s.DB.Query(`
	SELECT name, wins FROM players ORDER BY	wins DESC,name ASC;`)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}
	var players League
	for results.Next() {
		var player Player
		err := results.Scan(&player.Name, &player.Wins)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		players = append(players, player)
	}
	return players
}
