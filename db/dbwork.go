package db

import (
	"database/sql"
	"log"
	"sort"

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

func GetAllPlayers(db *sql.DB) []structs.Player {
	var result []structs.Player
	dbPlayers := FindAllPlayers(db)
	for _, player := range dbPlayers {
		result = append(result, structs.Player{Name: player.name, TgId: player.tgId})
	}
	return result
}

func GetMatchResult(db *sql.DB, matchId int64) *structs.Result {
	matchNullable := findOneMatch(db, matchId)
	if matchNullable == nil {
		return nil
	}

	matchTeamResults := findMatchTeamResultsByMatchId(db, matchId)
	matchPlayerRoles := findMatchPlayerRolesByMatchId(db, matchId)

	var playerIds []int64
	for _, matchPlayerRole := range matchPlayerRoles {
		playerIds = append(playerIds, matchPlayerRole.playerId)
	}
	players := findPlayersById(db, playerIds)
	playerByIdMap := make(map[int64]Player)
	for _, player := range players {
		playerByIdMap[player.id] = player
	}

	game := *structs.FindGameByName(matchNullable.game)

	playerRoles := make(map[structs.Player]string)
	for _, matchPlayerRole := range matchPlayerRoles {
		player := playerByIdMap[matchPlayerRole.playerId]
		playerRoles[structs.Player{Name: player.name, TgId: player.tgId}] = matchPlayerRole.role
	}

	sort.Slice(matchTeamResults, func(i int, j int) bool {
		return matchTeamResults[i].place < matchTeamResults[j].place
	})
	var teamOrder []string
	for _, matchTeamResult := range matchTeamResults {
		teamOrder = append(teamOrder, matchTeamResult.team)
	}

	return &structs.Result{Game: game, PlayerRoles: playerRoles, TeamOrder: teamOrder}
}

func InsertMatchResult(db *sql.DB, matchResult structs.Result) int64 {
	matchId := insertMatch(db, Match{id: -1, game: matchResult.Game.Name})
	log.Println("Inserted match with id", matchId)

	// TODO optimize, validate
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
