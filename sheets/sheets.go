package sheets

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	ts "github.com/brijs/teamsnap-team-gen/teamsnap"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Service struct {
	googleAppCredFile string
	tokFile           string
	spreadSheetID     string
	srv               *sheets.Service
}

func NewService() *Service {
	// ctx := context.Background()
	s := &Service{
		googleAppCredFile: "google_sheets_credentials.json",
		tokFile:           "token.json",
		spreadSheetID:     "1jJh3z_DrfJ-rktLmyXKjzhkm8K8oXXk8MZT9OL1xSM0",
	}

	b, err := ioutil.ReadFile(s.googleAppCredFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// only scoped to edit files created by Teamsnap-srca app
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.file")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	s.srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return s
}

func (s *Service) PublishMatch(nextMatch ts.Event, teamA []*ts.Player, teamB []*ts.Player, volunteers []*ts.Player, groupName string, teamAName string, teamBName string) {
	fmt.Println("PublishMatch")
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
	values = append(values, []interface{}{"Date", nextMatch.StartDate.Add(time.Minute * -15).Local().String()})
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
		values = append(values, []interface{}{p.FirstName + " " + p.LastName, i + 1, len(teamA) - i})
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
}

func (s *Service) GetPreferredTeam(teamName string, players []*ts.Player) {
	fmt.Println("GetPreferredTeamMappings")
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
		fmt.Println("No data found.")
	} else {
		// fmt.Println("Name, Name, Team")
		for _, row := range resp.Values {
			// Print columns A through C, which correspond to indices 0 and 4.
			if len(row) > 2 {
				// fmt.Printf("%s, %s, %s\n", row[0], row[1], row[2])
				key := row[0].(string) + " " + row[1].(string)
				if temp[key] != nil {
					temp[key].PreferredTeam = row[2].(string)
				}
			}
			if len(row) == 2 { // no mappings
				// fmt.Printf("%s, %s\n", row[0], row[1])
			}

		}
	}

	return
}

func (s *Service) GetTeamInfo(teamName string) (teamAName string, teamBName string) {
	fmt.Println("GetTeamInfo")
	ctx := context.Background()

	readRange := getTeamInfoRangeName(teamName)
	resp, err := s.srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// build a map for all players for quick lookup
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
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
