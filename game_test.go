package poker_test

import (
	"github.com/djangulo/go-fast"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server, _ := poker.NewPlayerServer(&poker.StubPlayerStore{}, poker.DummyGame)

		request := poker.NewGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		winner := "Ruth"
		srv, _ := poker.NewPlayerServer(store, poker.DummyGame)
		server := httptest.NewServer(srv)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		poker.AssertPlayerWin(t, store, winner)
	})
	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		game := &poker.GameSpy{}
		winner := "Ruth"
		server := httptest.NewServer(poker.MustMakePlayerServer(t, dummyPlayerStore, game))
		ws := poker.MustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		poker.WriteWSMessage(t, ws, "3")
		poker.WriteWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		poker.AssertGameStartedWith(t, game, 3)
		poker.AssertFinishCalledWith(t, game, winner)
	})
}
