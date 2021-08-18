package tgfunction

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	tg "github.com/brijs/teamsnap-team-gen/pkg/teamgen"
	log "github.com/sirupsen/logrus"
)

var (
	groupName          string    = "IntA"
	date               time.Time = time.Now().Add(time.Hour * -24)
	err                error
	opsNewSheet        bool
	teamRotationOffset int = -1
)

// For Google Cloud Function to work, this Function
// needs to be at the root level of the module

// /teamgen?groupName=IntA&afterDate=2021-08-01&rotationOffset=1
func TeamGen(w http.ResponseWriter, r *http.Request) {
	log.Info("TeamGen Function")

	log.Info("RawQuery=", r.URL.RawQuery)

	if rGroupName := r.URL.Query().Get("groupName"); rGroupName != "" {
		valid := false
		for _, g := range []string{"IntA", "IntB", "IntC", "IntD"} {
			if rGroupName == g {
				valid = true
			}
		}
		if !valid {
			log.Error("Invalid GroupName")
			http.Error(w, "Bad input groupName. Must be IntA|IntB|IntC|IntD", http.StatusInternalServerError)
			return
		}
		groupName = rGroupName
	}
	if rAfterDate := r.URL.Query().Get("afterDate"); rAfterDate != "" {
		layout := "2006-01-02"
		if date, err = time.Parse(layout, rAfterDate); err != nil {
			log.Error(err)
			if err != nil {
				http.Error(w, "Bad input date. Format yyyy-mm-dd:"+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	if rRotationOffset := r.URL.Query().Get("rotationOffset"); rRotationOffset != "" {
		if teamRotationOffset, err = strconv.Atoi(rRotationOffset); err != nil {
			log.Error(err)
			if err != nil {
				http.Error(w, "Bad input rotationOffset:"+err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}
	log.Infof("Params=> group=%s date=%v rotationOffset=%d", groupName, date, teamRotationOffset)

	spreadSheetID := tg.GenerateTeamsAndPublish(groupName, date, teamRotationOffset)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<p>Request successfully processed. Generated teams for <b>%s</b> group. </p><p>Click <a href=\"https://docs.google.com/spreadsheets/d/%s/edit#gid=660501010\">here</a> for the generated teams spreadsheet</p>", groupName, spreadSheetID)

}

// Test locally
// go run cmd/team-gen-function
// Browser=> localhost:8080/teamgen
//
// Deployed to GCP using
//  gcloud functions deploy TeamGen  --runtime go113 --trigger-http --allow-unauthenticated --project team-gen-function
