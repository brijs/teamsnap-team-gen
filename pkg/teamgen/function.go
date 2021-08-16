package tg

import (
	"fmt"
	"net/http"
	"strconv"
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

// /teamgen?groupName=IntA&afterDate=2021-08-01&rotationOffset=1
func TeamGen(w http.ResponseWriter, r *http.Request) {
	log.Info("TeamGen Function")

	if rGroupName := r.URL.Query().Get("groupName"); rGroupName != "" {
		valid := false
		for _, g := range []string{"IntA", "IntB", "IntC", "IntD"} {
			valid = rGroupName == g
		}
		if !valid {
			log.Error("Invalid GroupName")
		}
		groupName = rGroupName
	}
	if rAfterDate := r.URL.Query().Get("afterDate"); rAfterDate != "" {
		layout := "2006-01-02"
		if date, err = time.Parse(layout, rAfterDate); err != nil {
			log.Error(err)

		}
	}
	if rRotationOffset := r.URL.Query().Get("rotationOffset"); rRotationOffset != "" {
		if teamRotationOffset, err = strconv.Atoi(rRotationOffset); err != nil {
			log.Error(err)
		}
	}
	log.Infof("Params: group=%s date=%v rotationOffset=%d", groupName, date, teamRotationOffset)

	spreadSheetID := GenerateTeamsAndPublish(groupName, date, teamRotationOffset)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<p>Request successfully processed</p><p>Click <a href=\"https://docs.google.com/spreadsheets/d/%s/edit#gid=660501010\">here</a> for the generated teams spreadsheet</p>", spreadSheetID)

}
