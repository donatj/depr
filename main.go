package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mholt/archiver/v3"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	deprDir string
	deprLog *logfile

	set     = flag.String("set", "default", "Set to store to")
	msg     = flag.String("msg", "", "Description of items being stored")
	archive = flag.Bool("a", false, "archive the contents")
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

func init() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	valid := regexp.MustCompile(`^[\p{L}\d_\-]{1,63}$`)
	if !valid.MatchString(*set) {
		log.Fatal("invalid set name")
	}
}

type storeDetails struct {
	origPath string
	modPath  string
	archived bool
}

func main() {
	deprFiles := make(map[string]storeDetails)
	for _, e := range flag.Args() {
		stat, err := os.Stat(e)
		if err != nil {
			log.Fatal(err)
		}

		p, err := filepath.Abs(e)
		if err != nil {
			log.Fatal(err)
		}

		name := stat.Name()
		if *archive {
			zipfp := name + ".zip"
			err = archiver.Archive([]string{name}, zipfp)
			if err != nil {
				log.Fatal(err)
			}

			err = os.RemoveAll(name)
			if err != nil {
				log.Println(err)
				continue
			}

			deprFiles[p] = storeDetails{origPath: p, modPath: zipfp, archived: true}
		} else {
			deprFiles[p] = storeDetails{origPath: p, modPath: name}
		}
	}

	store(deprFiles, *set, *msg)
}
