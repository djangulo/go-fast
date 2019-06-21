package poker_test

import (
	"github.com/djangulo/go-fast"
	"github.com/djangulo/go-fast/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPostgreSQLPlayerStore integration test
func TestPostgreSQLStoreIntegration(t *testing.T) {
	store, removeStore := poker.NewPostgreSQLPlayerStore(
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
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
