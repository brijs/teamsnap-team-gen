package sheets

import (
	"bytes"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type formatConfig struct {
	SheetId        int
	TeamARowStart  int
	TeamBRowStart  int
	TeamBRowEnd    int
	TeamARowStart1 int
	TeamBRowStart1 int
}

func getJsonForFormatUpdates(cfg formatConfig) ([]byte, error) {
	req := `{
	"requests": [
		{
			"updateCells": {
				"range": {
					"sheetId": {{.SheetId}},
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
					"sheetId": {{.SheetId}},
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
						"horizontalAlignment": "LEFT",
						"textFormat": {
							"bold": true
						}
					}
				},
				"fields": "userEnteredFormat(horizontalAlignment,textFormat)",
				"range": {
					"sheetId": {{.SheetId}},
					"startRowIndex": 1,
					"endRowIndex": 4,
					"startColumnIndex": 1,
					"endColumnIndex": 2
				}
			}
		},
		{
			"repeatCell": {
				"cell": {
					"userEnteredFormat": {
						"backgroundColor": {
							"blue": 0.8,
							"green": 0.8,
							"red": 0.8
						}
					}
				},
				"fields": "userEnteredFormat(backgroundColor)",
				"range": {
					"sheetId": {{.SheetId}},
					"startRowIndex": 1,
					"endRowIndex": 6,
					"startColumnIndex": 0,
					"endColumnIndex": 3
				}
			}
		},
		{
			"repeatCell": {
				"cell": {
					"userEnteredFormat": {
						"backgroundColor": {
							"blue": 0.9,
							"green": 0.9,
							"red": 0.9
						},
						"textFormat": {
							"bold": true
						}
					}
				},
				"fields": "userEnteredFormat(backgroundColor,textFormat)",
				"range": {
					"sheetId": {{.SheetId}},
					"startRowIndex": {{.TeamARowStart}},
					"endRowIndex": {{.TeamARowStart1}},
					"startColumnIndex": 0,
					"endColumnIndex": 3
				}
			}
		},
		{
			"repeatCell": {
				"cell": {
					"userEnteredFormat": {
						"backgroundColor": {
							"blue": 0.9,
							"green": 0.9,
							"red": 0.9
						},
						"textFormat": {
							"bold": true
						}
					}
				},
				"fields": "userEnteredFormat(backgroundColor,textFormat)",
				"range": {
					"sheetId": {{.SheetId}},
					"startRowIndex": {{.TeamBRowStart}},
					"endRowIndex": {{.TeamBRowStart1}},
					"startColumnIndex": 0,
					"endColumnIndex": 3
				}
			}
		},
		{
			"updateBorders": {
				"range": {
					"sheetId": {{.SheetId}},
					"startRowIndex": {{.TeamARowStart}},
					"endRowIndex": {{.TeamBRowEnd}},
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
}	`

	t, err := template.New("sheet-format").Parse(req)
	if err != nil {
		log.Error("Erroring parsing sheet formating template", err)
		return nil, err
	}

	o := &bytes.Buffer{}
	if err = t.Execute(o, cfg); err != nil {
		log.Error("Error executing the sheet formating template", err)
		return nil, err
	}

	return o.Bytes(), nil
}
