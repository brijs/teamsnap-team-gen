package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	ts "github.com/brijs/teamsnap-team-gen/teamsnap"
)

var (
	teamNameLookup = map[string]int{
		"IntA": 6892639,
		"IntB": 6892639, // TODO
		"IntC": 6892639, // TODO
		"IntD": 6892639, // TODO
	}
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

func getRotation(e ts.Event, teamRotationOffset int) int {
	if teamRotationOffset >= 0 {
		return teamRotationOffset
	}

	// calculate based on number of weeks since Aug 8 (arbitrary)
	timeFormat := "2006-01-02"

	refDateStr := "2021-08-08"

	y, m, d := e.StartDate.Date()
	dStr := fmt.Sprintf("%04d-%02d-%02d", y, m, d)

	t, _ := time.Parse(timeFormat, dStr)
	f, _ := time.Parse(timeFormat, refDateStr)

	// fmt.Println(dStr, refDateStr, t, f)

	// count number of weeks since Aug 8
	duration := t.Sub(f)
	rotation := int(math.Ceil(float64(duration.Hours() / (24 * 7))))
	fmt.Println("Team Rotation offset = ", rotation)
	if rotation < 0 {
		return 0
	}
	return rotation
}

func printDebugInfo(players []*ts.Player) {
	fmt.Printf("Players len=%d\n", len(players))
	for _, p := range players {
		fmt.Printf("%+v\n", *p)
	}
}
