package poker

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
	game     Game
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store
	p.game = game
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(p.rootHandler))
	router.Handle("/game", http.HandlerFunc(p.serveGame))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router
	return p, nil
}

// PlayerStore interface
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

func (p *PlayerServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.GetLeague())
}
func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) getLeagueTable() League {
	return League{
		{Name: "Denis", Wins: 20},
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) serveGame(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, ioutil.Discard) //todo: Dont discard the blinds messages!

	_, winner, _ := conn.ReadMessage()
	p.game.Finish(string(winner))

}

// func GetPlayerScore(name string) int {
// 	if name == "Pepper" {
// 		return 20
// 	}
// 	if name == "Floyd" {
// 		return 10
// 	}
// 	return 0
// }
