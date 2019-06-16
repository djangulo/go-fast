package config

import (
	"os"
	fp "path/filepath"
)

var (
	// RootDir project root
	RootDir string
	// DatabaseFilename database string to pass into the store
	// TODO: adapt this so an sqlite or a postgres Connection String can be passed to the store
	DatabaseFilename string
)

func init() {
	pwd, _ := os.Getwd()
	RootDir = fp.Dir(fp.Dir(pwd))
	DatabaseFilename = fp.Join(RootDir, "game.db")
}
