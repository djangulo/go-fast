package main

import (
	"testing"
)

const (
	testDBName = ":memory:"
	name1      = "Peter"
	name2      = "John"
)

func TestCreatePlayer(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		store := NewSqlite3PlayerStore(testDBName)
		defer store.db.Close()

		player := Player{Name: name1, Score: 0}
		err := store.CreatePlayer(player)
		var p Player
		store.db.Where("name = ?", name1).First(&p)
		if p.Name != name1 {
			t.Errorf("got '%s' want '%s'", p.Name, name1)
		}
		assertNoError(t, err)
	})
	t.Run("error on unique field", func(t *testing.T) {
		store := NewSqlite3PlayerStore(testDBName)
		defer store.db.Close()

		player1 := Player{Name: name1, Score: 0}
		player2 := Player{Name: name1, Score: 0}
		player3 := Player{Name: name2, Score: 0}
		store.CreatePlayer(player1)
		store.CreatePlayer(player3)
		err := store.CreatePlayer(player2)
		assertError(t, err, ErrRecordAlreadyExists)
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
