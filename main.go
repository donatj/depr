package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

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

	deprLog, err = newLogfile(lf)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.Parse()

}

func main() {

	deprFiles := make(map[string]string)
	for _, e := range flag.Args() {
		stat, err := os.Stat(e)
		if err != nil {
			log.Fatal(err)
		}

		p, err := filepath.Abs(e)
		if err != nil {
			log.Fatal(err)
		}

		deprFiles[p] = stat.Name()
	}
	store(deprFiles)
}
