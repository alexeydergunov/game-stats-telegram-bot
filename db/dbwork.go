package db

import (
	"database/sql"
	"log"

	"example.com/structs"
)

func GetOrInsertPlayer(db *sql.DB, player structs.Player) Player {
	existingPlayer := FindOnePlayerByTgId(db, player.TgId)
	if existingPlayer != nil {
		log.Println("Found existing player", *existingPlayer)
		return *existingPlayer
	}
	playerId := insertPlayer(db, player)
	return *findOnePlayer(db, playerId)
}

func InsertMatchResult(db *sql.DB, matchResult structs.Result) int64 {
	matchId := insertMatch(db, Match{id: -1, game: matchResult.Game.Name})
	log.Println("Inserted match with id", matchId)

	// TODO optimize
	var matchPlayerRoles []MatchPlayerRole
	for player, role := range matchResult.PlayerRoles {
		playerId := FindOnePlayerByTgId(db, player.TgId).id
		matchPlayerRoles = append(matchPlayerRoles, MatchPlayerRole{id: -1, matchId: matchId, playerId: playerId, role: role})
	}
	for _, matchPlayerRole := range matchPlayerRoles {
		insertMatchPlayerRole(db, matchPlayerRole)
	}
	log.Println("Inserted", len(matchPlayerRoles), "match player roles for match", matchId)

	var matchTeamResults []MatchTeamResult
	for index, team := range matchResult.TeamOrder {
		matchTeamResults = append(matchTeamResults, MatchTeamResult{id: -1, matchId: matchId, team: team, place: int64(index) + 1})
	}
	for _, matchTeamResult := range matchTeamResults {
		insertMatchTeamResult(db, matchTeamResult)
	}
	log.Println("Inserted", len(matchTeamResults), "match team results for match", matchId)

	return matchId
}
