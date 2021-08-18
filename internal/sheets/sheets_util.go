package sheets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	googleAppCredFile string = "google_sheets_credentials.json"
	tokFile           string = "token.json"
	spreadsheetID     string = "1jJh3z_DrfJ-rktLmyXKjzhkm8K8oXXk8MZT9OL1xSM0"
)

// Retrieve a token, saves the token, then returns the generated client.
func getSheetsService() (*sheets.Service, error) {
	// set up Oauth config - either to refresh_grant or auth_code flow
	clientSecret, ok := os.LookupEnv("SHEETS_CLIENT_SECRET")
	if !ok || clientSecret == "" {
		log.Fatalf("SHEETS_CLIENT_SECRET not set or is empty")
	}

	config := &oauth2.Config{
		ClientID:     "701503350435-5duo582esn5mdncb6uvbp8i662m2f9c0.apps.googleusercontent.com",
		ClientSecret: clientSecret,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{sheets.DriveFileScope},
		Endpoint:     endpoints.Google,
	}

	// read cached token if present - the access_token may be stale, in which case, it'll be refreshed
	tok, err := tokenFromEnvOrFile(tokFile)
	// check if nothing is cached ie this is the first-time we are initiating Oauth (auth_code flow)
	if os.IsNotExist(err) {
		// auth code flow
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	} else if err != nil {
		log.Fatalf("Error trying to load a cached token: %v", err)
	}

	// Setup a "reuseTokenSource"
	ts := config.TokenSource(context.TODO(), tok)

	return sheets.NewService(context.TODO(), option.WithTokenSource(ts), option.WithScopes(sheets.DriveFileScope))
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	log.Info(("Starting Oauth2 3-legged flow to generate token.."))
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
func tokenFromEnvOrFile(file string) (*oauth2.Token, error) {
	var err error
	tok := &oauth2.Token{}

	// check env var first
	tokenJSON, ok := os.LookupEnv("SHEETS_TOKEN")
	if ok && tokenJSON != "" {
		log.Debug(("Trying to parse Sheets Token from env"))
		err = json.NewDecoder(strings.NewReader(tokenJSON)).Decode(tok)
	} else {
		log.Debug(("Trying to parse Sheets Token from file"))
		// read from file
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(tok)
	}

	if err != nil {
		log.Warn("Failed to load token from env or file", err)
		return nil, err
	}

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
	log.Info("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getPreferredTeamRangeName(teamName string) string {
	return teamName + "!" + teamName + "_PreferredTeam"
}

func getTeamInfoRangeName(teamName string) string {
	return teamName + "!" + teamName + "_Info"
}
