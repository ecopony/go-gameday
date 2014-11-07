package main

import (
//	_ "github.com/lib/pq"
	"log"
	"bytes"
	"os"
	"io/ioutil"
	"net/http"
	"encoding/xml"
//	"database/sql"
	s "strings"
)

type Epg struct {
	Date string `xml:"id,attr"`
	LastModified string `xml:"last_modified,attr"`
	DisplayTimeZone string `xml:"display_time_zone,attr"`
	EpgGames []EpgGame `xml:"game"`
}

/*
gids look like: gid_2014_07_22_nynmlb_seamlb_1

Doesn't yet handle doubleheader days. It'll just return the first match it finds for the team.
 */
func (e Epg) GidForTeam(teamCode string) string {
	for _, game := range e.EpgGames {
		if s.Contains(game.Gameday, s.Join([]string{"_", teamCode, "mlb_"}, "")) {
			return "gid_" + game.Gameday
		}
	}
	return "" // return an error here as well?
}

type EpgGame struct {
	CalendarEventId string `xml:"calendar_event_id,attr"`
	Start string `xml:"start,attr"`
	Id string `xml:"id,attr"`
	Gameday string `xml:"gameday,attr"`
}

type Game struct {
	XMLName xml.Name `xml:"game"`
	GameType string `xml:"type,attr"`
	LocalGameTime string `xml:"local_game_time,attr"`
	Teams []Team `xml:"team"`
	Stadium Stadium `xml:"stadium"`
}

type Team struct {
	XMLName xml.Name `xml:"team"`
	TeamType string `xml:"type,attr"`
	Code string `xml:"code,attr"`
	FileCode string `xml:"file_code,attr"`
}

type Stadium struct {
	XMLName xml.Name `xml:"stadium"`
	Id string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

func main() {
	args := os.Args[1:]
	if (len(args) != 2) {
		log.Fatal("Usage: gameday teamCode date")
	}

	teamCode := args[0]
	date := args[1]

	log.Println("Fetching game for " + teamCode + " on " + date)

	epgResp, err := http.Get(epgUrl(date))
	if err != nil {
		log.Fatal(err)
	}
	defer epgResp.Body.Close()
	epgBody, err := ioutil.ReadAll(epgResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var epg Epg
	xml.Unmarshal(epgBody, &epg)
	log.Println("Fetching from: " + gameUrl(date, epg.GidForTeam(teamCode)))

	resp, err := http.Get(gameUrl(date, epg.GidForTeam(teamCode)))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var game Game
	xml.Unmarshal(body, &game)
	log.Println(resp.Status)
	log.Println(string(body))

//	Assumes a pg database exists named go-gameday, a role that can access it.
//	Assumes a table called pitches with a character column called code.
//	db, err := sql.Open("postgres", "user=go-gameday dbname=go-gameday sslmode=disable")
//	issues := db.Ping()
//	log.Println(issue)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	rows, err :=  db.Query("SELECT code FROM pitches")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer rows.Close()
//
//	for rows.Next() {
//		var code string
//		err = rows.Scan(&code)
//		log.Println(code)
//	}
}

func baseUrl() string {
	return "http://gd2.mlb.com/components/game/mlb/"
}

func dateUrl(date string) string {
	var buffer bytes.Buffer
	buffer.WriteString(baseUrl())
	buffer.WriteString(datePath(date))
	return buffer.String()
}

func epgUrl(date string) string {
	var buffer bytes.Buffer
	buffer.WriteString(dateUrl(date))
	buffer.WriteString("/epg.xml")
	return buffer.String()
}

func gameDirectoryUrl(date string, gid string) string { // parse the date out of the gid to not have to pass both around
	var buffer bytes.Buffer
	buffer.WriteString(baseUrl())
	buffer.WriteString(datePath(date))
	buffer.WriteString("/")
	buffer.WriteString(gid)
	buffer.WriteString("/")
	return buffer.String()
}

func gameUrl(date string, gid string) string { // parse the date out of the gid to not have to pass both around
	return gameDirectoryUrl(date, gid) + "game.xml"
}

func datePath(date string) string {
	// firx this to be date parsing, validating
	datePieces := s.Split(date, "-")
	var buffer bytes.Buffer
	buffer.WriteString("year_")
	buffer.WriteString(datePieces[0])
	buffer.WriteString("/month_")
	buffer.WriteString(datePieces[1])
	buffer.WriteString("/day_")
	buffer.WriteString(datePieces[2])
	return buffer.String()
}

