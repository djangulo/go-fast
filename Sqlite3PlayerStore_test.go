package main

import (
	"testing"
)

const (
	testDBName = "test_db.db"
	name1      = "Peter"
	name2      = "John"
)

func TestCreatePlayer(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		store, removeStore := NewSqlite3PlayerStore(testDBName)
		defer removeStore()

		player := Player{Name: name1, Wins: 0}
		err := store.CreatePlayer(player)
		var p Player
		store.db.Where("name = ?", name1).First(&p)
		if p.Name != name1 {
			t.Errorf("got '%s' want '%s'", p.Name, name1)
		}
		assertNoError(t, err)
	})
	t.Run("error on creating existing name", func(t *testing.T) {
		store, removeStore := NewSqlite3PlayerStore(testDBName)
		defer removeStore()

		player1 := Player{Name: name1, Wins: 0}
		player2 := Player{Name: name1, Wins: 0}
		store.CreatePlayer(player1)
		err := store.CreatePlayer(player2)
		assertError(t, err, ErrRecordAlreadyExists)
	})
	t.Run("store wins for existing player", func(t *testing.T) {
		store, removeStore := NewSqlite3PlayerStore(testDBName)
		defer removeStore()

		player := Player{Name: name1, Wins: 0}
		store.CreatePlayer(player)
		store.RecordWin(name1)
		got := store.GetPlayerScore(name1)
		want := 1
		assertScoreEquals(t, got, want)
	})
	t.Run("store wins for new player", func(t *testing.T) {
		store, removeStore := NewSqlite3PlayerStore(testDBName)
		defer removeStore()

		store.RecordWin(name1)
		got := store.GetPlayerScore(name1)
		want := 1
		assertScoreEquals(t, got, want)
	})

}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but ditn't want one")
	}
}
