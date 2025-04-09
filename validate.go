package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
)

type validateCmd struct {
}

func (*validateCmd) Name() string     { return "validate" }
func (*validateCmd) Synopsis() string { return "Validate the backing database." }
func (*validateCmd) Usage() string {
	return `validate:
	  Validate the backing database.
  `
}

func (st *validateCmd) SetFlags(f *flag.FlagSet) {
}

func (st *validateCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	i, err := deprLog.Iter()
	if err != nil {
		log.Printf("error opening log file: %s", err)
		return subcommands.ExitFailure
	}

	line := 0
	log.Println("Validating log file...")
	for logEntry, err := range i {
		line++
		if err != nil {
			log.Printf("error reading log file: %s", err)
			return subcommands.ExitFailure
		}

		log.Println(logEntry.New)

		p := filepath.Join(deprDir, logEntry.New)
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			log.Printf("error: file '%s' does not exist on line %d", p, line)
			return subcommands.ExitFailure
		}

	}
	log.Println("Log file validated.")

	return subcommands.ExitSuccess
}
