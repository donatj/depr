package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func store(files map[string]string) {
	now := time.Now()

	for ff, f := range files {
		ffp, fp, err := getMvPath(now, f)
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

		deprLog.Append(ff, ffp, now)
	}
}

func getMvPath(now time.Time, name string) (string, string, error) {
	j := 0
	for {
		j++
		if j > 255 {
			return "", "", fmt.Errorf("Too many retries")
		}

		p := filepath.Join(
			now.Format("2006"),
			now.Format("01"),
			now.Format("02"),
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
