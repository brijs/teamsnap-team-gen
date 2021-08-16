package ts

import (
	"time"

	cj "github.com/brijs/teamsnap-team-gen/internal/collectionjson"
	log "github.com/sirupsen/logrus"
)

func mapToEvent(d []cj.DataType) Event {
	e := Event{}
	for _, v := range d {
		log.Trace(v.Name, "=>", v.Value)
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
				log.Error(err)
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
			log.Trace(v.Name, "=>", v.Value)
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
				log.Trace(v.Name, "%T %v=>\n", v.Value, v.Value)
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
			log.Trace(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "member_id":
				p.Id = uint64(v.Value.(float64))
				// p.Id = v.Value.(string)
			case "status_code":
				log.Trace(v.Name, "%T %v=>\n", v.Value, v.Value)
				p.IsAvailable = v.Value != nil && v.Value.(float64) != 0
			}
		}
		// update back in the players slice
		if temp[p.Id] != nil {
			temp[p.Id].IsAvailable = p.IsAvailable
		} else {
			log.Trace("WARN: couldn't update availability info for %+v\n", p)
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
			log.Trace(v.Name, "=>", v.Value)
			switch n := v.Name; n {
			case "member_id":
				if v.Value != nil {
					p.Id = uint64(v.Value.(float64))
					p.IsVolunteer = true
				}
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
			log.Warn("WARN: couldn't update assignment info for %+v\n", p)
		}

	}

}
