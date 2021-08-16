# teamsnap-team-gen

## Usage
```zsh

# Set the Teamsnap access token (contact repo owner)
$ export TEAMSNAP_TOKEN=xxx


# Print usage
$ ./teamsnap-team-gen --help

Usage of ./teamsnap-team-gen:
 Split available players for the specified team & date for an upcoming game

  -date value
    	Specify reference date (eg 2021/08/14). The script will find the first upcoming match after that date
  -groupname value
    	Specify one of the valid team names (IntA|IntB|IntC|IntD)
  -newSheet
    	Create a new Google Spreadsheet. (admin usage only)
  -rotateTeamOrder int
    	Enter a positive integer (optional) (default -1)


```

## Example run
```zsh
$ ./teamsnap-team-gen --groupname IntA

INFO[0000] Running for Teamsnap Team = (xxxx IntA), for date=2021-08-15 17:10:06.110296 -0700 PDT m=+0.019290270
INFO[0000] GetAllPlayersInTeam
INFO[0001] GetUpcomingEvent
INFO[0001] Event => {Id:xxxx TeamId:xxxx Location:GRMS - San Ramon Notes:Red ball, white jersey game. StartDate:2021-08-15 21:00:00 +0000 UTC Uniform:Whites LeagueName:Intermediate A Games}
INFO[0001] GetAvailability
INFO[0001] GetAssignments
INFO[0002] GetPreferredTeamMappings
INFO[0003] GetTeamInfo
INFO[0004] AssignTeamsToAvailablePlayers
INFO[0004] GetVolunteers
INFO[0004] PublishMatch
INFO[0005] Successfully completed generated teams for IntA
```

## Output

> - The generated teams are output to [this sheet](https://docs.google.com/spreadsheets/d/1jJh3z_DrfJ-rktLmyXKjzhkm8K8oXXk8MZT9OL1xSM0/edit#gid=2101538123)
