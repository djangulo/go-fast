package poker_test

import (
	"github.com/djangulo/go-fast"
	"reflect"
	"testing"
)

const (
	testDatabaseHost     = "localhost"
	testDatabasePort     = "5432"
	testDatabaseName     = "test_database"
	testDatabaseUser     = "postgres"
	testDatabasePassword = "abcd1234"
	name1                = "Peter"
	name2                = "John"
)

func TestPostgreSQLPlayerStore(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		store, removeStore := poker.NewPostgreSQLPlayerStore(
			testDatabaseHost,
			testDatabasePort,
			testDatabaseUser,
			testDatabaseName,
			testDatabasePassword,
		)
		defer removeStore()

		player := poker.Player{Name: name1, Wins: 0}
		err := store.CreatePlayer(player)
		var p poker.Player
		store.DB.Where("name = ?", name1).First(&p)
		if p.Name != name1 {
			t.Errorf("got '%s' want '%s'", p.Name, name1)
		}
		poker.AssertNoError(t, err)
	})
	t.Run("error on creating existing name", func(t *testing.T) {
		store, removeStore := poker.NewPostgreSQLPlayerStore(
			testDatabaseHost,
			testDatabasePort,
			testDatabaseUser,
			testDatabaseName,
			testDatabasePassword,
		)
		defer removeStore()

		player1 := poker.Player{Name: name1, Wins: 0}
		player2 := poker.Player{Name: name1, Wins: 0}
		store.CreatePlayer(player1)
		err := store.CreatePlayer(player2)
		poker.AssertError(t, err, poker.ErrRecordAlreadyExists)
	})
	t.Run("store wins for existing player", func(t *testing.T) {
		store, removeStore := poker.NewPostgreSQLPlayerStore(
			testDatabaseHost,
			testDatabasePort,
			testDatabaseUser,
			testDatabaseName,
			testDatabasePassword,
		)
		defer removeStore()

		player := poker.Player{Name: name1, Wins: 0}
		store.CreatePlayer(player)
		store.RecordWin(name1)
		got := store.GetPlayerScore(name1)
		want := 1
		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for new player", func(t *testing.T) {
		store, removeStore := poker.NewPostgreSQLPlayerStore(
			testDatabaseHost,
			testDatabasePort,
			testDatabaseUser,
			testDatabaseName,
			testDatabasePassword,
		)
		defer removeStore()

		store.RecordWin(name1)
		got := store.GetPlayerScore(name1)
		want := 1
		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("league should return ordered by wins, then name", func(t *testing.T) {
		store, removeStore := poker.NewPostgreSQLPlayerStore(
			testDatabaseHost,
			testDatabasePort,
			testDatabaseUser,
			testDatabaseName,
			testDatabasePassword,
		)
		defer removeStore()

		store.CreatePlayer(poker.Player{Name: "A", Wins: 1})
		store.CreatePlayer(poker.Player{Name: "B", Wins: 3})
		store.CreatePlayer(poker.Player{Name: "C", Wins: 2})
		store.CreatePlayer(poker.Player{Name: "D", Wins: 1})
		got := store.GetLeague()
		want := poker.League{
			{Name: "B", Wins: 3},
			{Name: "C", Wins: 2},
			{Name: "A", Wins: 1},
			{Name: "D", Wins: 1},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got \n%v \nwant \n%v", got, want)
		}
	})

}
