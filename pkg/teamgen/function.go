package tg

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	groupName          string    = "IntA"
	date               time.Time = time.Now().Add(time.Hour * -100)
	err                error
	opsNewSheet        bool
	teamRotationOffset int
)

// HelloWorld writes "Hello, World!" to the HTTP response.
func TeamGen(w http.ResponseWriter, r *http.Request) {
	log.Info("TeamGen Function")
	fmt.Fprint(w, "Hello, World!\n")
	fmt.Fprintf(w, "Click %s for the generated teams spreadsheet", "https://docs.google.com/spreadsheets/d/1jJh3z_DrfJ-rktLmyXKjzhkm8K8oXXk8MZT9OL1xSM0/edit#gid=2101538123")

	GenerateTeamsAndPublish(groupName, date, teamRotationOffset)
}
