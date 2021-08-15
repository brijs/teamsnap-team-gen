package ts

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	cj "github.com/brijs/teamsnap-team-gen/collectionjson"
)

// type Id uint64

type Player struct {
	Id                uint64 // string //uint64
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

func mapToEvent(d []cj.DataType) Event {
	e := Event{}
	for _, v := range d {
		// fmt.Println(v.Name, "=>", v.Value)
		switch n := v.Name; n {
		case "id":
			e.Id = uint64(v.Value.(float64))
		case "team_id":
			e.TeamId = uint64(v.Value.(float64))
		case "notes":
			if v.Value != nil {
				e.Notes = v.Value.(string)
			}
		case "uniform":
			e.Uniform = v.Value.(string)
		case "location_name":
			e.Location = v.Value.(string)
		case "start_date":
			// Monday, January 2, 2006 10:04:05 PM
			// 2021-08-15T21:00:00Z
			// 2021-08-15T21:00:00Z08:00 -- madeup
			// 2006-01-02T15:04:05Z07:00
			// t, err := time.Parse("2006-01-02T15:04:05Z", v.Value.(string))
			// todo: maybe use time_zone_offset
			t, err := time.Parse("2006-01-02T15:04:05Z07:00", v.Value.(string))

			if err != nil {
				fmt.Println(err)
			}
			e.StartDate = t
		case "opponent_name":
			e.LeagueName = v.Value.(string)
		}

	}
	return e
}

func mapToPlayers(i []cj.ItemType) []*Player {
	var players []*Player
	for _, d := range i {
		p, isPlayer := Player{}, true

		for _, v := range d.Data {
			// fmt.Println(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "id":
				p.Id = uint64(v.Value.(float64))
				// p.Id = v.Value.(string)
			case "first_name":
				p.FirstName = v.Value.(string)
			case "last_name":
				p.LastName = v.Value.(string)
			case "is_non_player":
				isPlayer = !v.Value.(bool)
			case "status_code":
				// fmt.Printf(v.Name, "%T %v=>\n", v.Value, v.Value)
				p.IsAvailable = v.Value != nil && v.Value.(float64) != 0
			}
		}
		if isPlayer {
			players = append(players, &p)
		}

	}
	return players
}

func mapAvailabilityToPlayers(i []cj.ItemType, players []*Player) {
	// build a map for all players for quick lookup
	var temp map[uint64]*Player = make(map[uint64]*Player)
	for _, p := range players {
		// p := p
		temp[p.Id] = p
	}

	for _, d := range i {
		p := Player{}
		for _, v := range d.Data {
			// fmt.Println(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "member_id":
				p.Id = uint64(v.Value.(float64))
				// p.Id = v.Value.(string)
			case "status_code":
				// fmt.Printf(v.Name, "%T %v=>\n", v.Value, v.Value)
				p.IsAvailable = v.Value != nil && v.Value.(float64) != 0
			}
		}
		// update back in the players slice
		if temp[p.Id] != nil {
			temp[p.Id].IsAvailable = p.IsAvailable
		} else {
			fmt.Printf("WARN: couldn't update availability info for %+v\n", p)
		}

	}

}
func mapAssignmentsToPlayers(i []cj.ItemType, players []*Player) {
	// build a map for all players for quick lookup
	var temp map[uint64]*Player = make(map[uint64]*Player)
	for _, p := range players {
		// p := p
		temp[p.Id] = p
	}

	for _, d := range i {
		p := Player{}
		for _, v := range d.Data {
			// fmt.Println(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "member_id":
				p.Id = uint64(v.Value.(float64))
				p.IsVolunteer = true
				// p.Id = v.Value.(string)
			case "position":
				p.VolunteerPosition = int(v.Value.(float64))
			case "description":
				p.VolunteerDesc = v.Value.(string)
			}
		}
		// update back in the players slice
		if temp[p.Id] != nil {
			tempP := temp[p.Id]
			tempP.IsVolunteer = p.IsVolunteer
			tempP.VolunteerDesc = p.VolunteerDesc
			tempP.VolunteerPosition = p.VolunteerPosition
		} else {
			fmt.Printf("WARN: couldn't update assignment info for %+v\n", p)
		}

	}

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

func (c *Client) GetUpcomingEvent2() (e Event, err error) {
	fmt.Println("GetUpcomingEvent")

	// test event
	req, err := http.NewRequest("GET", c.baseURL+"/events/246172835", nil)
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
	fmt.Printf("Event => %+v\n", e)

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
