package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

type logfile struct {
	filename string
}

func newLogfile(filename string) (*logfile, error) {
	if stat, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating", filename)
			file, err := os.Create(filename)
			defer file.Close()
			if err != nil {
				return nil, fmt.Errorf("error creating '%s': %s", deprDir, err)
			}
		} else if stat.IsDir() {
			return nil, fmt.Errorf("Error: '%s' is a directory.", filename)
		}
	}

	return &logfile{
		filename: filename,
	}, nil
}

func (l *logfile) Append(oldPath, newPath string, now time.Time) {
	logf, err := os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, 0600)
	defer logf.Close()
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(logf)
	defer func() {
		w.Flush()
		err := w.Error()
		if err != nil {
			log.Fatal(err)
		}
	}()

	w.Write([]string{
		newPath,
		oldPath,
		now.String(),
	})
}
