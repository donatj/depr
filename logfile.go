package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type logfile struct {
	deprDir  string
	filename string
}

func newLogfile(deprDir, filename string) (*logfile, error) {
	if stat, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating", filename)
			file, err := os.Create(filename)
			defer file.Close()
			if err != nil {
				return nil, fmt.Errorf("error creating '%s': %s", deprDir, err)
			}
		} else if stat.IsDir() {
			return nil, fmt.Errorf("error: '%s' is a directory", filename)
		}
	}

	return &logfile{
		deprDir:  deprDir,
		filename: filename,
	}, nil
}

type deprlog struct {
	New   string
	Old   string
	Descr string `json:",omitempty"`

	Archived bool

	Now time.Time
}

func (l *logfile) Append(d deprlog /* oldPath, newPath, descr string, now time.Time */) {
	logf, err := os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, 0600)
	defer logf.Close()
	if err != nil {
		log.Fatal(err)
	}

	rel, err := filepath.Rel(l.deprDir, d.New)
	if err == nil {
		d.New = rel
	} else {
		log.Printf("error: failed to get relative directory, see latest line of '%s'", l.filename)
	}

	w := json.NewEncoder(logf)
	w.Encode(d)
}
