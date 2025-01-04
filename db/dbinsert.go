package db

import (
	"database/sql"
	"log"

	"example.com/structs"
)

func insertPlayer(db *sql.DB, player structs.Player) int64 {
	log.Println("Inserting new player ...")
	statement, err := db.Prepare(`INSERT INTO player(name, tg_id) VALUES (?, ?)`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	insertResult, err := statement.Exec(player.Name, player.TgId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return id
}

func insertMatch(db *sql.DB, match Match) int64 {
	log.Println("Inserting new match ...")
	statement, err := db.Prepare(`INSERT INTO match(game) VALUES (?)`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	insertResult, err := statement.Exec(match.game)
	if err != nil {
		log.Fatalln(err.Error())
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return id
}

func insertMatchPlayerRole(db *sql.DB, matchPlayerRole MatchPlayerRole) int64 {
	log.Println("Inserting new match player role ...")
	statement, err := db.Prepare(`INSERT INTO match_player_role(match_id, player_id, role) VALUES (?, ?, ?)`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	insertResult, err := statement.Exec(matchPlayerRole.matchId, matchPlayerRole.playerId, matchPlayerRole.role)
	if err != nil {
		log.Fatalln(err.Error())
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return id
}

func insertMatchTeamResult(db *sql.DB, matchTeamResult MatchTeamResult) int64 {
	log.Println("Inserting new match team result ...")
	statement, err := db.Prepare(`INSERT INTO match_team_result(match_id, team, place) VALUES (?, ?, ?)`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	insertResult, err := statement.Exec(matchTeamResult.matchId, matchTeamResult.team, matchTeamResult.place)
	if err != nil {
		log.Fatalln(err.Error())
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return id
}
