package main

import (
	"bufio"
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

	Hostname string

	Now time.Time
}

func (l *logfile) Append(d deprlog) error {
	logf, err := os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer logf.Close()

	rel, err := filepath.Rel(l.deprDir, d.New)
	if err == nil {
		d.New = rel
	} else {
		log.Printf("error: failed to get relative directory, see latest line of '%s'", l.filename)
	}

	w := json.NewEncoder(logf)
	w.Encode(d)

	return nil
}

func (l *logfile) Read() ([]deprlog, error) {
	logf, err := os.Open(l.filename)
	if err != nil {
		return nil, err
	}
	defer logf.Close()

	scanner := bufio.NewScanner(logf)
	var logs []deprlog
	for scanner.Scan() {
		var d deprlog
		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			return nil, err
		}
		logs = append(logs, d)
	}

	return logs, err
}
