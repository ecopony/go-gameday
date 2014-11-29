package gamedayapi

import (
	"bytes"
	"log"
	"os/user"
	s "strings"
)

const (
	// GamedayHostname is the hostname of the MLB gameday site
	GamedayHostname = "http://gd2.mlb.com"

	// GamedayBaseURL is the base URL of the MLB gameday files
	GamedayBaseURL = "http://gd2.mlb.com/components/game/mlb"
)

func datePath(date string) string {
	// firx this to be date parsing, validating
	datePieces := s.Split(date, "-")
	var buffer bytes.Buffer
	buffer.WriteString("/year_")
	buffer.WriteString(datePieces[0])
	buffer.WriteString("/month_")
	buffer.WriteString(datePieces[1])
	buffer.WriteString("/day_")
	buffer.WriteString(datePieces[2])
	return buffer.String()
}

func dateURL(date string) string {
	var buffer bytes.Buffer
	buffer.WriteString(GamedayBaseURL)
	buffer.WriteString(datePath(date))
	return buffer.String()
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func baseCachePath() string {
	return homeDir() + "/go-gameday-cache"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
