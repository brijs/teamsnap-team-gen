package sheets

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func CreateNewSheet() string {
	ctx := context.Background()

	b, err := ioutil.ReadFile("google_sheets_credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Scopes Explained
	//
	// Ideal setup
	// ------------
	//  Use https://www.googleapis.com/auth/drive.file
	//  This means that this app will only have access to files it creates. Least risk
	//  For this to work, you'd first need to Create file through this app.
	//  Files created by user directly will not be accessible
	//
	//  Using https://www.googleapis.com/auth/spreadsheets.readonly is broader than
	//  necessary, as it gives this app read access to all spreadsheets in Drive.
	//
	// NOTE 1 - When you change the scope in code, you need to delete token.json
	//  and go through OAuth flow again to regenerate it. Otherwise, the app will
	//  continue to use the old token.json
	//
	// NOTE 2 - in GCP, you also need to specify scope that this app CAN request.
	//   The app is identified by Client ID (& secret) & Project ID in GCP
	//   From the app, you can only request a subset of those scopes.

	deleteTokenFile()

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.file")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Create file - using JSON (Method 1)
	createData := []byte(`
	{
		"sheets": [
		  {
			"properties": {
			  "title": "Sample Tab"
			}
		  }
		],
		"properties": {
		  "title": "Create Spreadsheet using Sheets API v4"
		}
	  }
	`)
	var createReq sheets.Spreadsheet
	if err = json.Unmarshal(createData, &createReq); err != nil {
		log.Fatal(err)
	}

	// Method 2 - write the Request Struct  - more verbose
	// createReq := &sheets.Spreadsheet{
	// 	Properties: &sheets.SpreadsheetProperties{
	// 		Title: "teamsnap-srca-created2",
	// 	},
	// 	Sheets: []*sheets.Sheet{
	// 		&sheets.Sheet{
	// 			Properties: &sheets.SheetProperties{
	// 				Title: "boom",
	// 			},
	// 		},
	// 	},
	// 	// TODO: Add desired fields of the request body.

	// }

	resp, err := srv.Spreadsheets.Create(&createReq).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%#v\n", resp)

	return resp.SpreadsheetUrl
}
