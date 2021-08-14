package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func getTeamSnapToken() string {
	val, ok := os.LookupEnv("TEAMSNAP_TOKEN")
	if !ok {
		log.Fatalln("TEAMSNAP_TOKEN is not set")
		return val
	} else {
		return val
	}
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

var Usage = func() {
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n Split available players for the specified team & date for an upcoming game\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// flags
	var (
		teamName string    = "IntA"
		date     time.Time = time.Now()
		err      error
	)

	enumFlag(&teamName, "team", []string{"IntA", "IntB", "IntC", "IntD"}, "Specify one of the valid team names (IntA|IntB|IntC|IntD)\n")
	flag.Func("date", "Specify reference date (eg 2021/08/14). The script will find the first upcoming match after that date", func(flagValue string) error {
		layout := "2006/01/02"
		if date, err = time.Parse(layout, flagValue); err != nil {
			fmt.Println(err)
			return err
		}
		return err
	})
	flag.Usage = Usage

	flag.Parse()
	fmt.Println(teamName, date)

	// tsClient := ts.NewClient(getTeamSnapToken())

	// e, _ := tsClient.GetUpcomingEvent()
	// fmt.Printf("Event => %+v\n", e)

	// tsClient.GetAvailability(e)
}
