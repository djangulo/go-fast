package main

import (
	"flag"
	"fmt"
	"github.com/djangulo/go-fast"
	"github.com/djangulo/go-fast/config"
	"log"
	"time"
	"os"
	"net/http"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

var port string
var postgresTimeout int

func init() {
	const (
		defaultPort = "9000"
		defaultTimeout = 60
		portUsage   = "port to serve the app on, default '" + defaultPort + "'"
		timeoutUsage   = "time limit to wait for a postgres connection, default " + string(defaultPort)
	)
	flag.StringVar(&port, "port", defaultPort, portUsage)
	flag.StringVar(&port, "p", defaultPort, portUsage+" (shorthand)")
	flag.IntVar(&postgresTimeout, "postgres-timeout", defaultTimeout, timeoutUsage)
	flag.IntVar(&postgresTimeout, "t", defaultTimeout, timeoutUsage+" (shorthand)")
}

func postgresHealthcheck() <-chan error {
	c := make(chan error)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
    )
	go func() {
		for {
			db, err := gorm.Open("postgres", connStr)
			if err != nil {
				c <- err
			} else {
				if erro := db.DB().Ping(); erro != nil {
					c <- erro
				} else {
					c <- nil
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return c
}

func waitForDB() {
	c := postgresHealthcheck()
	timeout := time.After(time.Duration(postgresTimeout) * time.Second)
	for {
		select {
		case err := <-c:
			if err != nil {
				fmt.Printf("Waiting for postgres (connection error: %v)\n", err)
			} else {
				fmt.Println("Postgres is ready")
				return
			}
		case <-timeout:
			fmt.Println("Postgres timeout exceeded, aborting")
			os.Exit(1)
		}
	}
}


func main() {
	flag.Parse()
	waitForDB()
	fmt.Println("Listening at 127.0.0.1:" + port)
	store, removeStore := poker.NewPostgreSQLPlayerStore(
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
	)
	defer removeStore()
	server := poker.NewPlayerServer(store)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
