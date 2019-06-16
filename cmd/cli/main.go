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

	store, _ := poker.NewSqlite3PlayerStore(config.DatabaseFilename)
	defer store.DB.Close()

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()

}
