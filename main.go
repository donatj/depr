package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	deprDir string
	deprLog *logfile
)

func init() {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	deprDir = filepath.Join(dir, ".depr")

	if stat, err := os.Stat(deprDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(deprDir, 0755)
			if err != nil {
				log.Fatalf("Error creating '%s': %s", deprDir, err)
			}
		} else if !stat.IsDir() {
			log.Fatalf("Error: '%s' is not a directory.", deprDir)
		}
	}

	lf := filepath.Join(deprDir, "depr.log")

	deprLog, err = newLogfile(deprDir, lf)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "help")
	subcommands.Register(subcommands.FlagsCommand(), "help")
	subcommands.Register(subcommands.CommandsCommand(), "help")
	store := &storeCmd{}
	subcommands.Register(store, "")
	subcommands.Register(subcommands.Alias("s", store), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
