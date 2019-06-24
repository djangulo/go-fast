package main

import (
	"fmt"
	"github.com/djangulo/go-fast"
	"github.com/djangulo/go-fast/config"
	"os"
)

func main() {
	fmt.Println("Lets play poker")
	fmt.Println("Type {Name} wins to recond a win")

	store, removeStore := poker.NewPostgreSQLPlayerStore(
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
	)
	defer removeStore()

	game := poker.NewTexasHoldem(
		poker.BlindAlerterFunc(poker.StdOutAlerter),
		store,
	)

	poker.NewCLI(
		os.Stdin,
		os.Stdout,
		game,
	).PlayPoker()

}
