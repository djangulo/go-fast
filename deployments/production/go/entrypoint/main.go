package entrypoint

import (
	"flag"
	"fmt"
	"github.com/djangulo/go-fast/config"
	"time"
	"os"
	"os/exec"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

var postgresTimeout int

func init() {
	const (
		defaultTimeout = 60
		timeoutUsage   = "time limit to wait for a postgres connection, default " + string(defaultTimeout)
	)
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
	// pipe it forward
	fmt.Println(os.Args)
	exec.Command("go", os.Args[1:]...)
}
