package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var timeoutSeconds int

func init() {
	const (
		defaultTimeout = 60
		timeoutUsage   = "Seconds to wait for the PostgreSQL database, default " + string(defaultTimeout)
	)
	flag.IntVar(&timeoutSeconds, "t", defaultTimeout, timeoutUsage)
}

func postgresReady() <-chan int {
	c := make(chan int)
	go func() {
		connStr := fmt.Sprintf(
			"user=%s dbname=%s host=%s port=%s password=%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DB"),
		)
		for {
			_, err := sql.Open("postgres", connStr)
			if err != nil {
				c <- 0
			} else {
				c <- 1
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return c
}

func main() {
	flag.Parse()
	c := postgresReady()
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	for {
		select {
		case exitStatus := <-c:
			if exitStatus == 0 {
				fmt.Println("PostgreSQL is available...")
				os.Exit(0)
			} else {
				fmt.Println("Waiting for PostgreSQL to become available...")
			}
		case <-timeout:
			fmt.Printf("Timeout (%d seconds) exceeded, aborting...\n", timeoutSeconds)
			os.Exit(1)
		}
	}
}
