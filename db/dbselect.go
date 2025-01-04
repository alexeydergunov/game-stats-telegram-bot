package db

import (
	"database/sql"
	"log"
	"strings"

	"example.com/structs"
)

func findPlayersByName(db *sql.DB, players []structs.Player) []Player {
	log.Println("Finding", len(players), "players")
	var playerNames []any
	var questions []string
	for _, player := range players {
		playerNames = append(playerNames, player.Name)
		questions = append(questions, "?")
	}
	whereList := "(" + strings.Join(questions, ",") + ")"
	log.Println("whereList: " + whereList)
	statement, err := db.Prepare("SELECT * FROM player WHERE name in " + whereList) // TODO use tg_id
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(playerNames...)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Player
	for row.Next() {
		var player Player
		row.Scan(&player.id, &player.name, &player.tgId)
		log.Println("Found player:", player)
		result = append(result, player)
	}
	log.Println("Found", len(result), "players in DB")
	return result
}

func findOnePlayerByName(db *sql.DB, player structs.Player) *Player {
	log.Println("Finding player", player)
	result := findPlayersByName(db, []structs.Player{player})
	if len(result) == 0 {
		return nil
	}
	if len(result) == 1 {
		return &result[0]
	}
	log.Fatalln("Found", len(result), "players in DB instead of one")
	return nil
}

func findOnePlayer(db *sql.DB, playerId int64) *Player {
	log.Println("Finding player by id", playerId)
	statement, err := db.Prepare("SELECT * FROM player WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(playerId)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Player
	for row.Next() {
		var player Player
		row.Scan(&player.id, &player.name, &player.tgId)
		log.Println("Found player:", player)
		result = append(result, player)
	}
	if len(result) == 0 {
		return nil
	}
	if len(result) == 1 {
		return &result[0]
	}
	log.Fatalln("Found", len(result), "players in DB instead of one")
	return nil
}

func findOneMatch(db *sql.DB, matchId int64) *Match {
	log.Println("Finding match by id", matchId)
	statement, err := db.Prepare("SELECT * FROM match WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(matchId)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Match
	for row.Next() {
		var match Match
		row.Scan(&match.id, &match.game)
		log.Println("Found match:", match)
		result = append(result, match)
	}
	if len(result) == 0 {
		return nil
	}
	if len(result) == 1 {
		return &result[0]
	}
	log.Fatalln("Found", len(result), "matches in DB instead of one")
	return nil
}
