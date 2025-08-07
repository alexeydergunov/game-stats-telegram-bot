package db

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

func intToString(x int64) string {
	return strconv.FormatInt(x, 10)
}

func findPlayersById(db *sql.DB, ids []int64) []Player {
	log.Println("Finding", len(ids), "players")
	var idsStr []any
	var questions []string
	for _, id := range ids {
		idsStr = append(idsStr, intToString(id))
		questions = append(questions, "?")
	}
	whereList := "(" + strings.Join(questions, ",") + ")"
	log.Println("whereList: " + whereList)
	statement, err := db.Prepare("SELECT * FROM player WHERE id in " + whereList)
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(idsStr...)
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

func findPlayersByTgId(db *sql.DB, tgIds []int64) []Player {
	log.Println("Finding", len(tgIds), "players")
	var tgIdsStr []any
	var questions []string
	for _, tgId := range tgIds {
		tgIdsStr = append(tgIdsStr, intToString(tgId))
		questions = append(questions, "?")
	}
	whereList := "(" + strings.Join(questions, ",") + ")"
	log.Println("whereList: " + whereList)
	statement, err := db.Prepare("SELECT * FROM player WHERE tg_id in " + whereList)
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(tgIdsStr...)
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

func FindOnePlayerByTgId(db *sql.DB, tgId int64) *Player {
	log.Println("Finding player with tgId", tgId)
	result := findPlayersByTgId(db, []int64{tgId})
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

func findMatchTeamResultsByMatchId(db *sql.DB, matchId int64) []MatchTeamResult {
	log.Println("Finding match team results with matchId", matchId)
	statement, err := db.Prepare("SELECT * FROM match_team_result WHERE match_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(matchId)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []MatchTeamResult
	for row.Next() {
		var matchTeamResult MatchTeamResult
		row.Scan(&matchTeamResult.id, &matchTeamResult.matchId, &matchTeamResult.team, &matchTeamResult.place)
		log.Println("Found match team result:", matchTeamResult)
		result = append(result, matchTeamResult)
	}
	log.Println("Found", len(result), "match team results in DB")
	return result
}

func findMatchPlayerRolesByMatchId(db *sql.DB, matchId int64) []MatchPlayerRole {
	log.Println("Finding match player roles with matchId", matchId)
	statement, err := db.Prepare("SELECT * FROM match_player_role WHERE match_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	row, err := statement.Query(matchId)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []MatchPlayerRole
	for row.Next() {
		var matchPlayerRole MatchPlayerRole
		row.Scan(&matchPlayerRole.id, &matchPlayerRole.matchId, &matchPlayerRole.playerId, &matchPlayerRole.role)
		log.Println("Found match player role:", matchPlayerRole)
		result = append(result, matchPlayerRole)
	}
	log.Println("Found", len(result), "match team results in DB")
	return result
}
