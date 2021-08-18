package tg

import (
	"fmt"
	"math"
	"os"
	"time"

	sheets "github.com/brijs/teamsnap-team-gen/internal/sheets"
	ts "github.com/brijs/teamsnap-team-gen/internal/teamsnap"
	teamsplit "github.com/brijs/teamsnap-team-gen/internal/teamsplit"
	log "github.com/sirupsen/logrus"
)

var (
	teamNameLookup = map[string]int{
		"IntA": 6892639,
		"IntB": 6892640,
		"IntC": 7368380,
		"IntD": 7427307,
	}
)

func GenerateTeamsAndPublish(groupName string, date time.Time, teamRotationOffset int) string {
	teamId := teamNameLookup[groupName]
	log.Info("Running for Teamsnap Team = (", teamId, groupName, "), for date=", date)

	tsClient := ts.NewClient(getTeamSnapToken())

	// 1. Get  all players in team
	players, _ := tsClient.GetAllPlayersInTeam(teamId)
	printDebugInfo(players)

	// 2. Get Upcoming match
	nextMatch, _ := tsClient.GetUpcomingEvent(teamId, date)
	log.Infof("Event => %+v", nextMatch)

	// 3. Get Player availability
	tsClient.GetAvailability(nextMatch.Id, players)
	printDebugInfo(players)

	// 4. Get Volunteer assignments
	tsClient.GetAssignments(nextMatch.Id, teamId, players)
	printDebugInfo(players)

	sheetsService := sheets.NewService()
	// 5. Get Stick team pref
	sheetsService.GetPreferredTeam(groupName, players)
	teamAName, teamBName := sheetsService.GetTeamInfo(groupName)
	printDebugInfo(players)

	// 6. Split into teams
	teamA, teamB := teamsplit.AssignTeamsToAvailablePlayers(players, getRotation(nextMatch, teamRotationOffset), teamAName, teamBName)
	printDebugInfo(teamA)
	printDebugInfo(teamB)

	// 7. Get Volunteers
	volunteers := teamsplit.GetVolunteers(players)
	printDebugInfo(volunteers)

	// 8. Format / publish to spreadsheet
	sheetsService.PublishMatch(nextMatch, teamA, teamB, volunteers, groupName, teamAName, teamBName)

	log.Info("Successfully completed generated teams for ", groupName)
	return sheetsService.SpreadSheetID
}

func getTeamSnapToken() string {
	val, ok := os.LookupEnv("TEAMSNAP_TOKEN")
	if !ok {
		log.Fatalln("TEAMSNAP_TOKEN is not set")
		return val
	} else {
		return val
	}
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

	log.Trace(dStr, refDateStr, t, f)

	// count number of weeks since Aug 8
	duration := t.Sub(f)
	rotation := int(math.Ceil(float64(duration.Hours() / (24 * 7))))
	log.Debug("Team Rotation offset = ", rotation)
	if rotation < 0 {
		return 0
	}
	return rotation
}

func printDebugInfo(players []*ts.Player) {
	log.Debugf("Players len=%d\n", len(players))
	for _, p := range players {
		log.Debugf("%+v\n", *p)
	}
}
