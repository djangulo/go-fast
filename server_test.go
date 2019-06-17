package poker_test

import (
	"github.com/djangulo/go-fast"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETRoot(t *testing.T) {
	store := poker.NewStubPlayerStore(map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}, nil, nil)
	server := poker.NewPlayerServer(store)
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	poker.AssertStatus(t, response.Code, http.StatusOK)
	poker.AssertResponseBody(t, response.Body.String(), "Hello, World!")

}

func TestGETPlayers(t *testing.T) {
	store := poker.NewStubPlayerStore(map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}, nil, nil)
	server := poker.NewPlayerServer(store)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})
	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.NewStubPlayerStore(map[string]int{}, nil, nil)
	server := poker.NewPlayerServer(store)
	t.Run("it records wins when POSTed", func(t *testing.T) {
		player := "Pepper"
		request := poker.NewPostWinRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response.Code, http.StatusAccepted)
		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin, want %d", len(store.WinCalls), 1)
		}
		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner got '%s' want '%s'", store.WinCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := poker.League{
			{Name: "Jazmin", Wins: 4},
			{Name: "David", Wins: 3},
			{Name: "Elena", Wins: 6},
		}

		store := poker.NewStubPlayerStore(nil, nil, wantedLeague)
		server := poker.NewPlayerServer(store)

		request := poker.NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := poker.GetLeagueFromResponse(t, response.Body)
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
		poker.AssertContentType(t, response, poker.JsonContentType)
	})
}

// TestPostgreSQLPlayerStore integration test
func TestPostgreSQLStoreIntegration(t *testing.T) {
	store, removeStore := poker.NewPostgreSQLPlayerStore(
		testDatabaseHost,
		testDatabasePort,
		testDatabaseUser,
		testDatabaseName,
		testDatabasePassword,
	)
	defer removeStore()

	server := poker.NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewGetScoreRequest(player))
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "3")
	})
	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewLeagueRequest())
		poker.AssertStatus(t, response.Code, http.StatusOK)
		got := poker.GetLeagueFromResponse(t, response.Body)
		want := poker.League{
			{Name: "Pepper", Wins: 3},
		}
		poker.AssertLeague(t, got, want)
	})
}
