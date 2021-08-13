package main

import (
	"fmt"
	"log"
	"os"

	ts "github.com/brijs/teamsnap-team-gen/teamsnap"
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

func main() {

	tsClient := ts.NewClient(getTeamSnapToken())

	e, _ := tsClient.GetUpcomingEvent()
	fmt.Printf("Event => %+v\n", e)

	tsClient.GetAvailability(e)
}
