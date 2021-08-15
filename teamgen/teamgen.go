package teamgen

import (
	"fmt"
	"math"
	"sort"
	"strings"

	ts "github.com/brijs/teamsnap-team-gen/teamsnap"
)

type adj struct {
	a2b int
	b2a int
	n2a int
	n2b int
}

func calculateAdjustments(total, a, b int) adj {
	n := total - (a + b)
	// fmt.Println("*** ", total, ":", a, b)
	a2b, b2a, n2a, n2b := 0, 0, 0, 0
	a_tgt := int(math.Ceil(float64(total) / 2))
	b_tgt := int(math.Floor(float64(total) / 2))
	//fmt.Println(a_tgt, b_tgt)

	if a > a_tgt {
		a2b = a - a_tgt
		n2b = n
	} else if b > b_tgt {
		b2a = b - b_tgt
		n2a = n
	}

	if a <= a_tgt && b <= b_tgt {
		n2a = a_tgt - a
		n2b = b_tgt - b
	}

	return adj{a2b, b2a, n2a, n2b}
}

func filterByTeamNameAndSort(players []*ts.Player, teamName string, rotation int) (ret []*ts.Player) {
	// filter Available players
	for _, p := range players {
		if p.IsAvailable && p.PreferredTeam == teamName {
			ret = append(ret, p)
		}
	}
	if len(ret) == 0 {
		return
	}

	// sort in descending order
	sort.Slice(ret, func(i, j int) bool {
		return strings.Compare(ret[i].FirstName+ret[i].LastName, ret[j].FirstName+ret[j].LastName) == 1
	})

	// rotate
	rotation = rotation % len(ret)
	rotation = len(ret) - rotation // rotate forward; 1st=>2nd, 2nd=>3rd
	ret = append(ret[rotation:], ret[0:rotation]...)

	return
}
func AssignTeamsToAvailablePlayers(players []*ts.Player, rotation int, teamAName string, teamBName string) (teamA []*ts.Player, teamB []*ts.Player) {
	fmt.Println("AssignTeamsToAvailablePlayers")
	total, a, b := 0, 0, 0

	// get current team counts
	for _, p := range players {
		if p.IsAvailable {
			if p.PreferredTeam == teamAName {
				a, total = a+1, total+1
			} else if p.PreferredTeam == teamBName {
				b, total = b+1, total+1
			} else {
				total = total + 1
			}
		}
	}

	// calculate adjustments
	adj := calculateAdjustments(total, a, b)

	// do the team adjustments
	for _, p := range players {
		if p.IsAvailable {
			if p.PreferredTeam == teamAName {
				if adj.a2b > 0 {
					p.PreferredTeam = teamBName
					adj.a2b = adj.a2b - 1
				}
			} else if p.PreferredTeam == teamBName {
				if adj.b2a > 0 {
					p.PreferredTeam = teamAName
					adj.b2a = adj.b2a - 1
				}

			} else {
				if adj.n2a > 0 {
					p.PreferredTeam = teamAName
					adj.n2a = adj.n2a - 1
				} else if adj.n2b > 0 {
					p.PreferredTeam = teamBName
					adj.n2b = adj.n2b - 1
				}
			}
		}
	}

	teamA = filterByTeamNameAndSort(players, teamAName, rotation)
	teamB = filterByTeamNameAndSort(players, teamBName, rotation)
	return
}

func GetVolunteers(players []*ts.Player) (ret []*ts.Player) {
	fmt.Println("GetVolunteers")

	// filter Available players
	for _, p := range players {
		if p.IsVolunteer {
			ret = append(ret, p)
		}
	}
	if len(ret) == 0 {
		return
	}

	// sort in descending order
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].VolunteerPosition < ret[j].VolunteerPosition
	})
	return
}
