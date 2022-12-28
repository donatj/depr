package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/google/subcommands"
	"github.com/mholt/archiver"
)

type storeCmd struct {
	set     string
	msg     string
	archive bool
}

func (*storeCmd) Name() string     { return "store" }
func (*storeCmd) Synopsis() string { return "Store the given files." }
func (*storeCmd) Usage() string {
	return `store [-a] [-msg="message"] [<files>...]:
	Store the given files.
  `
}

func (st *storeCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&st.set, "set", "default", "Set to store to")
	f.StringVar(&st.msg, "msg", "", "Description of items being stored")
	f.BoolVar(&st.archive, "a", false, "archive the contents")
}

func (st *storeCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	valid := regexp.MustCompile(`^[\p{L}\d_\-]{1,63}$`)
	if !valid.MatchString(st.set) {
		log.Println("invalid set name")
		return subcommands.ExitFailure
	}

	deprFiles := make(map[string]storeDetails)
	for _, e := range f.Args() {
		stat, err := os.Stat(e)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}

		p, err := filepath.Abs(e)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}

		name := stat.Name()
		if st.archive {
			zipfp := name + ".zip"
			err = archiver.Archive([]string{name}, zipfp)
			if err != nil {
				log.Println(err)
				return subcommands.ExitFailure
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

	store(deprFiles, st.set, st.msg)

	return subcommands.ExitSuccess
}

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

		deprLog.Append(deprlog{
			Old:   f.origPath,
			New:   ffp,
			Descr: descr,
			Now:   now,

			Archived: f.archived,
		})
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

type storeDetails struct {
	origPath string
	modPath  string
	archived bool
}
