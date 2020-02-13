package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mholt/archiver/v3"
)

func store(files map[string]string, set, descr string) {
	now := time.Now()

	for ff, f := range files {
		ffp, fp, err := getMvPath(now, f, set)
		_ = fp
		if err != nil {
			log.Println(err)
			continue
		}

		err = os.MkdirAll(fp, 0755)
		if err != nil {
			log.Println(err)
			continue
		}

		err = os.Rename(ff, ffp)
		if err != nil {
			log.Println(err)
			continue
		}

		if *archive {
			zipfp := ffp + ".zip"

			err = archiver.Archive([]string{ffp}, zipfp)
			if err != nil {
				log.Println(err)
				log.Printf("leaving unarchived file: %s", ffp)
				deprLog.Append(ff, ffp, descr, now)
				continue
			}

			deprLog.Append(ff, zipfp, descr, now)

			err = os.RemoveAll(ffp)
			if err != nil {
				log.Println(err)
				continue
			}

			continue
		}

		deprLog.Append(ff, ffp, descr, now)
	}
}

func getMvPath(now time.Time, name, set string) (string, string, error) {
	j := 0
	for {
		j++
		if j > 255 {
			return "", "", fmt.Errorf("Too many retries")
		}

		p := filepath.Join(
			set,
			now.Format("2006-01-02"),
			strconv.Itoa(j),
		)

		fp := filepath.Join(deprDir, p)
		ffp := filepath.Join(fp, name)

		if _, err := os.Stat(ffp); err != nil {
			if os.IsNotExist(err) {
				return ffp, fp, nil
			}
		}
	}
}
