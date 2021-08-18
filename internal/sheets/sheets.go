package sheets

import (
	"encoding/json"
	"fmt"
	"time"

	ts "github.com/brijs/teamsnap-team-gen/internal/teamsnap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

type Service struct {
	googleAppCredFile string
	tokFile           string
	SpreadSheetID     string
	srv               *sheets.Service
}

func NewService() *Service {
	// ctx := context.Background()
	s := &Service{
		googleAppCredFile: "google_sheets_credentials.json",
		tokFile:           "token.json",
		SpreadSheetID:     "1jJh3z_DrfJ-rktLmyXKjzhkm8K8oXXk8MZT9OL1xSM0",
	}
	var err error
	s.srv, err = getSheetsService()
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return s
}

func (s *Service) PublishMatch(nextMatch ts.Event, teamA []*ts.Player, teamB []*ts.Player, volunteers []*ts.Player, groupName string, teamAName string, teamBName string) {
	log.Info("PublishMatch")
	ctx := context.Background()

	rangeData := groupName + "_Match!A1:Z1000"
	values := [][]interface{}{}

	// First Clear Data
	rb1 := &sheets.ClearValuesRequest{}
	_, err := s.srv.Spreadsheets.Values.Clear(spreadsheetID, rangeData, rb1).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	// Match
	values = append(values, []interface{}{nextMatch.LeagueName})
	values = append(values, []interface{}{"Location", nextMatch.Location})
	// hardcode 15 min earlier than startDate
	loc, _ := time.LoadLocation("America/Los_Angeles")
	sd := nextMatch.StartDate.Add(time.Minute * -15).In(loc)
	values = append(values, []interface{}{"Date", sd.Format("Jan 2, 2006")})
	values = append(values, []interface{}{"Reporting Time", sd.Format(time.Kitchen)})
	values = append(values, []interface{}{"Uniform", nextMatch.Uniform})
	values = append(values, []interface{}{"Note", nextMatch.Notes})

	values = append(values, []interface{}{""})
	values = append(values, []interface{}{""})

	// Volunteer
	for _, v := range volunteers {
		values = append(values, []interface{}{v.VolunteerDesc, v.FirstName + " " + v.LastName})
	}
	values = append(values, []interface{}{""})
	values = append(values, []interface{}{""})

	// team A
	values = append(values, []interface{}{"Team " + teamAName, "Batting Order", "Bowling Order"})
	for i, p := range teamA {
		values = append(values, []interface{}{p.FirstName + " " + p.LastName, i + 1, len(teamA) - i})
	}
	values = append(values, []interface{}{""})

	// team B
	values = append(values, []interface{}{"Team " + teamBName, "Batting Order", "Bowling Order"})
	for i, p := range teamB {
		values = append(values, []interface{}{p.FirstName + " " + p.LastName, i + 1, len(teamB) - i})
	}
	values = append(values, []interface{}{""})

	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: values,
	})
	_, err = s.srv.Spreadsheets.Values.BatchUpdate(spreadsheetID, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// Format the sheets
	var sheetId int // TODO: dynamically resolve these
	switch groupName {
	case "IntA":
		sheetId = 660501010
	case "IntB":
		sheetId = 1959274600
	case "IntC":
		sheetId = 257356858
	case "IntD":
		sheetId = 1437245619
	}
	teamsRowStart := 5 + 2 + 2 + 1 + len(volunteers)
	teamsRowEnd := teamsRowStart + len(teamA) + len(teamB) + 3
	// clear all formatting
	// background for top row
	// Bold Location/Date/Time
	// Borders around teams
	r1 := []byte(fmt.Sprintf(`
	{
		"requests": [
			{
				"updateCells": {
					"range": {
						"sheetId": %d,
						"startRowIndex": 5
					},
					"fields": "userEnteredFormat"
				}
			},
			{
				"repeatCell": {
					"cell": {
						"userEnteredFormat": {
							"backgroundColor": {
								"blue": 0.2,
								"green": 0.2,
								"red": 0.2
							},
							"horizontalAlignment": "CENTER",
							"textFormat": {
								"foregroundColor": {
									"red": 1.0,
									"green": 1.0,
									"blue": 1.0
								},
								"bold": true
							}
						}
					},
					"fields": "userEnteredFormat(backgroundColor,horizontalAlignment,textFormat)",
					"range": {
						"sheetId": %d,
						"startRowIndex": 0,
						"endRowIndex": 1,
						"startColumnIndex": 0,
						"endColumnIndex": 1
					}
				}
			},
			{
				"repeatCell": {
					"cell": {
						"userEnteredFormat": {
							"textFormat": {
								"bold": true
							}
						}
					},
					"fields": "userEnteredFormat(textFormat)",
					"range": {
						"sheetId": %d,
						"startRowIndex": 1,
						"endRowIndex": 4,
						"startColumnIndex": 1,
						"endColumnIndex": 2
					}
				}
			},
			{
				"updateBorders": {
					"range": {
						"sheetId": %d,
						"startRowIndex": %d,
						"endRowIndex": %d,
						"startColumnIndex": 0,
						"endColumnIndex": 3
					},
					"top": {
						"style": "SOLID",
						"width": 1,
						"color": {
							"black": 1.0
						}
					},
					"bottom": {
						"style": "SOLID",
						"width": 1,
						"color": {
							"black": 1.0
						}
					},
					"left": {
						"style": "SOLID",
						"width": 1,
						"color": {
							"black": 1.0
						}
					},
					"right": {
						"style": "SOLID",
						"width": 1,
						"color": {
							"black": 1.0
						}
					},
					"innerHorizontal": {
						"style": "SOLID",
						"width": 1,
						"color": {
							"black": 1.0
						}
					},
					"innerVertical": {
						"style": "SOLID",
						"width": 1,
						"color": {
							"black": 1.0
						}
					},
					"fields": "*"
				}
			}
		]
	}	`, sheetId, sheetId, sheetId, sheetId, teamsRowStart, teamsRowEnd))
	rr1 := &sheets.BatchUpdateSpreadsheetRequest{}
	if err := json.Unmarshal(r1, rr1); err != nil {
		log.Error(err)
	}
	_, err = s.srv.Spreadsheets.BatchUpdate(s.SpreadSheetID, rr1).Context(ctx).Do()
	if err != nil {
		log.Error(err)
	}

}

func (s *Service) GetPreferredTeam(teamName string, players []*ts.Player) {
	log.Info("GetPreferredTeamMappings")
	ctx := context.Background()

	readRange := getPreferredTeamRangeName(teamName)
	resp, err := s.srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// build a map for all players for quick lookup
	var temp map[string]*ts.Player = make(map[string]*ts.Player)
	for _, p := range players {
		// p := p
		temp[p.FirstName+" "+p.LastName] = p
	}

	if len(resp.Values) == 0 {
		log.Warn("No data found.")
	} else {
		// fmt.Println("Name, Name, Team")
		for _, row := range resp.Values {
			// Print columns A through C, which correspond to indices 0 and 4.
			if len(row) > 2 {
				log.Trace("%s, %s, %s\n", row[0], row[1], row[2])
				key := row[0].(string) + " " + row[1].(string)
				if temp[key] != nil {
					temp[key].PreferredTeam = row[2].(string)
				}
			}
			if len(row) == 2 { // no mappings
				log.Trace("%s, %s\n", row[0], row[1])
			}

		}
	}

	return
}

func (s *Service) GetTeamInfo(teamName string) (teamAName string, teamBName string) {
	log.Info("GetTeamInfo")
	ctx := context.Background()

	readRange := getTeamInfoRangeName(teamName)
	resp, err := s.srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// build a map for all players for quick lookup
	if len(resp.Values) == 0 {
		log.Warn("No data found.")
	}
	if len(resp.Values) < 3 || len(resp.Values[1]) < 2 || len(resp.Values[2]) < 2 {
		log.Fatalln("TeamInfo Range not found or missing rows.")
	}

	return resp.Values[1][1].(string), resp.Values[2][1].(string)

}

// General Sheets API Usage Example
// -----------------------------------
//
// 1. Construct RequestBody
// rb = &sheets.BatchUpdateValuesRequest{} // rb: requestBody
//
// 2. Append Data to rb.Data
// rb.Data = append(rb.Data, &sheets.ValueRange{})
//
// 3. Call Method
// sheetsService.Spreadsheets.Values.BatchUpdate(spreadsheetID, rb).Context(ctx).Do()
//
// Common Classes
//  - sheets.XXXRequest{}    	  // ReqBody structs
//  - Spreadsheets.Values.XXX()   // methods
//  - sheets.ValueRange			  // Common util struct
//
// Methods (in sheetsService.Spreadsheets.Values) & Structs (in sheets)
//   - BatchGet			|
//   - BatchUpdate		| BatchUpdateValuesRequest
//   - Clear			| ClearValuesRequest
//   -

// Note: Each Method has it's corresponding RequesBody Struct
