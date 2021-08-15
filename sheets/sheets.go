package sheets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	ts "github.com/brijs/teamsnap-team-gen/teamsnap"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	tokFile       string = "token.json"
	spreadsheetID string = "1jJh3z_DrfJ-rktLmyXKjzhkm8K8oXXk8MZT9OL1xSM0"
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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Delete token file if it exists
func deleteTokenFile() {
	if err := os.Remove(tokFile); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error deleting token file", err)
	}
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func PersistToSheetTest(players []*ts.Player) {
	rangeData := "Test!A1:Z1000" // 1000 rows
	values := [][]interface{}{
		{"Id", "FirstName", "LastName", "Team", "Av", "Vol", "Desc", "VolPos"}}
	for _, p := range players {
		values = append(values, []interface{}{p.Id, p.FirstName, p.LastName, p.PreferredTeam, p.IsAvailable, p.IsVolunteer, p.VolunteerDesc, p.VolunteerPosition})
	}

	persistToSheet(rangeData, values)
	log.Println("Sheets: Done updating Batting Stats.")
}

func persistToSheet(rangeData string, values [][]interface{}) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("google_sheets_credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// only scoped to edit files created by Teamsnap-srca app

	// config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.file")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: values,
	})
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetID, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(ctx)
	return
}

func (s *Service) GetPreferredTeam(players []*ts.Player) {
	ctx := context.Background()

	readRange := "IntA_PreferredTeam"
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
		fmt.Println("Name, Name, Team")
		for _, row := range resp.Values {
			// Print columns A through C, which correspond to indices 0 and 4.
			if len(row) > 2 {
				fmt.Printf("%s, %s, %s\n", row[0], row[1], row[2])
				key := row[0].(string) + " " + row[1].(string)
				if temp[key] != nil {
					temp[key].PreferredTeam = row[2].(string)
				}
			}
			if len(row) == 2 {
				fmt.Printf("%s, %s\n", row[0], row[1])
			}

		}
	}

	return
}

func (s *Service) PublishMatch(nextMatch ts.Event, teamA []*ts.Player, teamB []*ts.Player, volunteers []*ts.Player) {
	ctx := context.Background()

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	rangeData := "NextMatch!A1:Z1000"
	values := [][]interface{}{}

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
	values = append(values, []interface{}{"Team Avengers", "Batting Order", "Bowling Order"})
	for i, p := range teamA {
		values = append(values, []interface{}{p.FirstName + " " + p.LastName, i + 1, len(teamA) - i})
	}
	values = append(values, []interface{}{""})

	// team B
	values = append(values, []interface{}{"Team Defenders", "Batting Order", "Bowling Order"})
	for i, p := range teamB {
		values = append(values, []interface{}{p.FirstName + " " + p.LastName, i + 1, len(teamA) - i})
	}
	values = append(values, []interface{}{""})

	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: values,
	})
	_, err := s.srv.Spreadsheets.Values.BatchUpdate(spreadsheetID, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
}
