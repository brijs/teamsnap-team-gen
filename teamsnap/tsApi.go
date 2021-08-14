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
	Id            uint64 // string //uint64
	FirstName     string
	LastName      string
	PreferredTeam string
	Available     bool
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
			e.Notes = v.Value.(string)
		case "uniform":
			e.Uniform = v.Value.(string)
		case "":
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

func mapToPlayers(i []cj.ItemType) []Player {
	var players []Player
	for _, d := range i {
		p := Player{}
		for _, v := range d.Data {
			// fmt.Println(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "id":
				p.Id = uint64(v.Value.(float64))
				// p.Id = v.Value.(string)
			case "first_name":
				p.FirstName = v.Value.(string)
			case "last_name":
				p.FirstName = v.Value.(string)

			case "status_code":
				// fmt.Printf(v.Name, "%T %v=>\n", v.Value, v.Value)
				p.Available = v.Value != nil && v.Value.(float64) != 0
			}
		}
		players = append(players, p)

	}
	return players
}

func mapToMembers(i []cj.ItemType) []Player {
	var players []Player
	for _, d := range i {
		p := Player{}
		for _, v := range d.Data {
			// fmt.Println(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "id":
				p.Id = uint64(v.Value.(float64))
				// p.Id = v.Value.(string)
			case "status_code":
				// fmt.Printf(v.Name, "%T %v=>\n", v.Value, v.Value)
				p.Available = v.Value != nil && v.Value.(float64) != 0
			}
		}
		players = append(players, p)

	}
	return players
}

func (c *Client) GetAllPlayersInTeam(teamId int) (players []Player, err error) {
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

func (c *Client) GetUpcomingEvent() (e Event, err error) {
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

func (c *Client) GetAvailability(e Event) (err error) {
	fmt.Println("GetMembers")

	// Create a new request using http
	req, err := http.NewRequest("GET", c.baseURL+"availabilities/search?event_id="+strconv.FormatUint(e.Id, 10), nil)
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

	players := mapToMembers(res.Collection.Items)
	fmt.Printf("Players => %+v\n", players)

	return err
}
