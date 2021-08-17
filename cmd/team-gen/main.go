package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	sheets "github.com/brijs/teamsnap-team-gen/internal/sheets"
	tg "github.com/brijs/teamsnap-team-gen/pkg/teamgen"
	log "github.com/sirupsen/logrus"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n Split available players for the specified team & date for an upcoming game\n\n", os.Args[0])
	flag.PrintDefaults()
}

func enumFlag(target *string, name string, safelist []string, usage string) {
	flag.Func(name, usage, func(flagValue string) error {
		for _, allowedValue := range safelist {
			if flagValue == allowedValue {
				*target = flagValue
				return nil
			}
		}

		return fmt.Errorf("must be one of %v", safelist)
	})
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	// log.SetLevel(log.DebugLevel)

	// log.SetReportCaller(true)
}

func main() {
	// flags
	var (
		groupName          string    = "IntA"
		date               time.Time = time.Now()
		err                error
		opsNewSheet        bool
		logLevel           string
		teamRotationOffset int
	)

	enumFlag(&groupName, "groupname", []string{"IntA", "IntB", "IntC", "IntD"}, "Specify one of the valid team names (IntA|IntB|IntC|IntD)")
	flag.Func("date", "Specify reference date (eg 2021/08/14). The script will find the first upcoming match after that date", func(flagValue string) error {
		layout := "2006/01/02"
		if date, err = time.Parse(layout, flagValue); err != nil {
			log.Error(err)
			return err
		}
		return err
	})
	flag.IntVar(&teamRotationOffset, "rotateTeamOrder", -1, "Enter a positive integer (optional)")
	flag.BoolVar(&opsNewSheet, "newSheet", false, "Create a new Google Spreadsheet. (admin usage only)")
	flag.StringVar(&logLevel, "logLevel", "info", "fatal error info debug trace")

	flag.Usage = Usage

	flag.Parse()

	if opsNewSheet {
		log.Info("Creating a new sheet & exiting")
		url := sheets.CreateNewSheet()
		log.Info("New Spreadsheet URL: ", url)
		return
	}
	if logLevel != "" {
		l, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Fatalf("Invalid log level", err)
		}
		log.Info("Setting Log Level: ", logLevel)
		log.SetLevel(l)
	}

	tg.GenerateTeamsAndPublish(groupName, date, teamRotationOffset)
	log.Info("Main done")
}
