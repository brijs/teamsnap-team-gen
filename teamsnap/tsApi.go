package ts

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// type Id uint64

type Player struct {
	Id                uint64
	FirstName         string
	LastName          string
	PreferredTeam     string
	IsAvailable       bool
	IsVolunteer       bool
	VolunteerDesc     string
	VolunteerPosition int
}

type Event struct {
	Id         uint64
	TeamId     uint64
	Location   string
	Notes      string
	StartDate  time.Time
	Uniform    string
	LeagueName string
}

func (c *Client) GetAllPlayersInTeam(teamId int) (players []*Player, err error) {
	fmt.Println("GetAllPlayersInTeam")

	// test event
	req, err := http.NewRequest("GET", c.baseURL+fmt.Sprintf("members/search?team_id=%d", teamId), nil)
	if err != nil {
		log.Fatalln("Error creating Request.\n[ERROR] -", err)
	}

	res, err := c.sendRequest(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}

	if len(res.Collection.Items) == 0 {
		return players, fmt.Errorf("Event was not found")
	}

	players = mapToPlayers(res.Collection.Items)
	return
}

func (c *Client) GetUpcomingEvent(teamId int, date time.Time) (e Event, err error) {
	fmt.Println("GetUpcomingEvent")

	log.Println(c.baseURL + fmt.Sprintf("/events/search?started_after=%s&page_size=1&team_id=%d", date.Format("2006-01-02T15:04"), teamId))

	req, err := http.NewRequest("GET", c.baseURL+fmt.Sprintf("/events/search?started_after=%s&page_size=1&team_id=%d", date.Format("2006-01-02T15:04"), teamId), nil)
	if err != nil {
		log.Fatalln("Error creating Request.\n[ERROR] -", err)
	}

	res, err := c.sendRequest(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}

	if len(res.Collection.Items) == 0 {
		return e, fmt.Errorf("Event was not found")
	}
	e = mapToEvent(res.Collection.Items[0].Data)

	return e, err
}

func (c *Client) GetAvailability(eventId uint64, players []*Player) (err error) {
	fmt.Println("GetAvailability")

	// Create a new request using http
	req, err := http.NewRequest("GET", c.baseURL+"availabilities/search?event_id="+strconv.FormatUint(uint64(eventId), 10), nil)
	if err != nil {
		log.Fatalln("Error creating Request.\n[ERROR] -", err)
	}

	res, err := c.sendRequest(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}

	if len(res.Collection.Items) == 0 {
		return fmt.Errorf("No Players were not found")
	}

	mapAvailabilityToPlayers(res.Collection.Items, players)

	return err
}

// https://api.teamsnap.com/v3/assignments/search?event_id=246172835&team_id=6892639
func (c *Client) GetAssignments(eventId uint64, teamId int, players []*Player) (err error) {
	fmt.Println("GetAssignments")

	// Create a new request using http
	req, err := http.NewRequest("GET", fmt.Sprintf(c.baseURL+"assignments/search?event_id=%d&team_id=%d", eventId, teamId), nil)
	if err != nil {
		log.Fatalln("Error creating Request.\n[ERROR] -", err)
	}

	res, err := c.sendRequest(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}

	if len(res.Collection.Items) == 0 {
		return fmt.Errorf("No Players were not found")
	}

	mapAssignmentsToPlayers(res.Collection.Items, players)

	return err
}
