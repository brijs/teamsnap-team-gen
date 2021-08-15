package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	sheets "github.com/brijs/teamsnap-team-gen/sheets"
	teamgen "github.com/brijs/teamsnap-team-gen/teamgen"
	ts "github.com/brijs/teamsnap-team-gen/teamsnap"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n Split available players for the specified team & date for an upcoming game\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// flags
	var (
		groupName          string    = "IntA"
		date               time.Time = time.Now()
		err                error
		opsNewSheet        bool
		teamRotationOffset int
	)

	enumFlag(&groupName, "group", []string{"IntA", "IntB", "IntC", "IntD"}, "Specify one of the valid team names (IntA|IntB|IntC|IntD)\n")
	flag.Func("date", "Specify reference date (eg 2021/08/14). The script will find the first upcoming match after that date", func(flagValue string) error {
		layout := "2006/01/02"
		if date, err = time.Parse(layout, flagValue); err != nil {
			fmt.Println(err)
			return err
		}
		return err
	})
	flag.BoolVar(&opsNewSheet, "newSheet", false, "Create a new Google Spreadsheet. (admin usage only)")
	flag.IntVar(&teamRotationOffset, "rotateTeamOrder", -1, "Enter a positive integer (optional)")

	flag.Usage = Usage

	flag.Parse()

	teamId := teamNameLookup[groupName]
	fmt.Println("Running for Teamsnap Team = (", teamId, groupName, "), for date=", date)

	if opsNewSheet {
		log.Println("Creating a new sheet & exiting")
		url := sheets.CreateNewSheet()
		log.Println("New Spreadsheet URL: ", url)
		return
	}

	tsClient := ts.NewClient(getTeamSnapToken())

	// 1. Get  all players in team
	players, _ := tsClient.GetAllPlayersInTeam(teamId)
	printDebugInfo(players)

	// 2. Get Upcoming match
	nextMatch, _ := tsClient.GetUpcomingEvent(teamId, date)
	fmt.Printf("Event => %+v\n", nextMatch)

	// 3. Get Player availability
	tsClient.GetAvailability(nextMatch.Id, players)
	// printDebugInfo(players)

	// 4. Get Volunteer assignments
	tsClient.GetAssignments(nextMatch.Id, teamId, players)
	// printDebugInfo(players)

	sheetsService := sheets.NewService()
	// 5. Get Stick team pref
	sheetsService.GetPreferredTeam(groupName, players)
	teamAName, teamBName := sheetsService.GetTeamInfo(groupName)
	// printDebugInfo(players)

	// 6. Split into teams
	teamA, teamB := teamgen.AssignTeamsToAvailablePlayers(players, getRotation(nextMatch, teamRotationOffset), teamAName, teamBName)
	printDebugInfo(teamA)
	printDebugInfo(teamB)

	// 7. Get Volunteers
	volunteers := teamgen.GetVolunteers(players)
	printDebugInfo(volunteers)

	// 8. Format / publish to spreadsheet
	sheetsService.PublishMatch(nextMatch, teamA, teamB, volunteers, groupName, teamAName, teamBName)

}
