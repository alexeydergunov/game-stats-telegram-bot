package db

import (
	"database/sql"
	"log"
)

func FindAllPlayers(db *sql.DB) []Player {
	log.Println("Finding all players")
	statement, err := db.Prepare("SELECT * FROM player")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Player
	for row.Next() {
		var obj Player
		row.Scan(&obj.id, &obj.name, &obj.tgId)
		log.Println("Found player:", obj)
		result = append(result, obj)
	}
	log.Println("Found", len(result), "players in DB")
	return result
}

func FindAllMatches(db *sql.DB) []Match {
	log.Println("Finding all matches")
	statement, err := db.Prepare("SELECT * FROM match")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Match
	for row.Next() {
		var obj Match
		row.Scan(&obj.id, &obj.game)
		log.Println("Found match:", obj)
		result = append(result, obj)
	}
	log.Println("Found", len(result), "matches in DB")
	return result
}

func FindAllMatchPlayerRoles(db *sql.DB) []MatchPlayerRole {
	log.Println("Finding all match player roles")
	statement, err := db.Prepare("SELECT * FROM match_player_role")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []MatchPlayerRole
	for row.Next() {
		var obj MatchPlayerRole
		row.Scan(&obj.id, &obj.matchId, &obj.playerId, &obj.role)
		log.Println("Found match player role:", obj)
		result = append(result, obj)
	}
	log.Println("Found", len(result), "match player roles in DB")
	return result
}

func FindAllMatchTeamResults(db *sql.DB) []MatchTeamResult {
	log.Println("Finding all match team results")
	statement, err := db.Prepare("SELECT * FROM match_team_result")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []MatchTeamResult
	for row.Next() {
		var obj MatchTeamResult
		row.Scan(&obj.id, &obj.matchId, &obj.team, &obj.place)
		log.Println("Found match team result:", obj)
		result = append(result, obj)
	}
	log.Println("Found", len(result), "match team results in DB")
	return result
}
