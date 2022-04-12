package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func store(files map[string]storeDetails, set, descr string) {
	now := time.Now()

	for _, f := range files {
		ffp, fp, err := getMvPath(now, f.modPath, set)
		if err != nil {
			log.Println(err)
			continue
		}

		err = os.MkdirAll(fp, 0755)
		if err != nil {
			log.Println(err)
			continue
		}

		err = os.Rename(f.modPath, ffp)
		if err != nil {
			log.Println(err)
			continue
		}

		deprLog.Append(f.origPath, ffp, descr, now)
	}
}

func getMvPath(now time.Time, name, set string) (string, string, error) {
	j := 0
	for {
		j++
		if j > 255 {
			return "", "", fmt.Errorf("too many retries - %d", j)
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
