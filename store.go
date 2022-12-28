package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/chzyer/readline"
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

		fullPath, err := filepath.Abs(e)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}

		shortName := stat.Name()
		if st.archive {
			tmpdir, err := os.MkdirTemp("", "depr-")
			if err != nil {
				log.Println("error creating temporary directory", err)
				return subcommands.ExitFailure
			}
			// todo - clean up old temporary directory

			zipPath := path.Join(tmpdir, shortName+".zip")
			err = archiver.Archive([]string{fullPath}, zipPath)
			if err != nil {
				log.Printf("error archiving '%s': %s", fullPath, err)
				return subcommands.ExitFailure
			}

			// todo - move to POST-success
			err = os.RemoveAll(fullPath)
			if err != nil {
				log.Println(err)
				continue
			}

			deprFiles[fullPath] = storeDetails{origPath: fullPath, modPath: zipPath, archived: true}
		} else {
			deprFiles[fullPath] = storeDetails{origPath: fullPath, modPath: fullPath}
		}
	}

	msg := st.msg
	if msg == "" {
		rl, err := readline.NewEx(&readline.Config{
			Prompt: "description (enter for none): ",

			DisableAutoSaveHistory: true,
		})
		if err != nil {
			log.Println("error opening readline", err)
			return subcommands.ExitFailure
		}
		defer rl.Close()

		rs := rl.Line()
		if rs.Error != nil {
			log.Println(rs.Error)
			return subcommands.ExitFailure
		}

		msg = rs.Line
	}

	store(deprFiles, st.set, msg)

	return subcommands.ExitSuccess
}

func store(files map[string]storeDetails, set, descr string) {
	now := time.Now()

	for _, f := range files {
		ffp, fp, err := getMvPath(now, f.modPath, set)
		if err != nil {
			log.Printf("path generation error: %s", err)
			continue
		}

		err = os.MkdirAll(fp, 0755)
		if err != nil {
			log.Printf("error making storage path: %s", err)
			continue
		}

		err = os.Rename(f.modPath, ffp)
		if err != nil {
			log.Printf("error moving file: %s", err)
			continue
		}

		host, err := os.Hostname()
		if err != nil {
			log.Printf("error getting hostname: %s -- storing as 'unknown'", err)
			host = "unknown"
		}

		err = deprLog.Append(deprlog{
			Old:   f.origPath,
			New:   ffp,
			Descr: descr,
			Now:   now,

			Archived: f.archived,

			Hostname: host,
		})
		if err != nil {
			log.Printf("error writing log file: %s", err)
			continue
		}
	}
}

func getMvPath(now time.Time, fullPath, set string) (string, string, error) {

	stat, err := os.Stat(fullPath)
	if err != nil {
		return "", "", err
	}

	shortName := stat.Name()

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
		ffp := filepath.Join(fp, shortName)

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
