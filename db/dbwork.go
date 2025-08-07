package db

import (
	"database/sql"
	"log"
	"slices"
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

func makeResult(game structs.Game, playerByIdMap map[int64]Player, matchTeamResults []MatchTeamResult, matchPlayerRoles []MatchPlayerRole) structs.Result {
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

	return structs.Result{Game: game, PlayerRoles: playerRoles, TeamOrder: teamOrder}
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

	result := makeResult(game, playerByIdMap, matchTeamResults, matchPlayerRoles)
	return &result
}

func GetMatchResultsByGame(db *sql.DB, gameName string) map[int64]structs.Result {
	game := *structs.FindGameByName(gameName)

	allMatches := FindAllMatches(db)
	var goodMatchIds []int64
	for _, match := range allMatches {
		if match.game == game.Name {
			goodMatchIds = append(goodMatchIds, match.id)
		}
	}
	slices.Sort(goodMatchIds)

	matchTeamResults := FindAllMatchTeamResults(db)
	matchPlayerRoles := FindAllMatchPlayerRoles(db)

	var playerIds []int64
	var matchTeamResultsByMatchId = make(map[int64][]MatchTeamResult)
	var matchPlayerRolesByMatchId = make(map[int64][]MatchPlayerRole)
	for _, matchId := range goodMatchIds {
		matchTeamResultsByMatchId[matchId] = []MatchTeamResult{}
		matchPlayerRolesByMatchId[matchId] = []MatchPlayerRole{}
	}
	for _, matchTeamResult := range matchTeamResults {
		matchId := matchTeamResult.matchId
		_, contains := slices.BinarySearch(goodMatchIds, matchId)
		if contains {
			matchTeamResultsByMatchId[matchId] = append(matchTeamResultsByMatchId[matchId], matchTeamResult)
		}
	}
	for _, matchPlayerRole := range matchPlayerRoles {
		matchId := matchPlayerRole.matchId
		_, contains := slices.BinarySearch(goodMatchIds, matchId)
		if contains {
			matchPlayerRolesByMatchId[matchId] = append(matchPlayerRolesByMatchId[matchId], matchPlayerRole)
			playerIds = append(playerIds, matchPlayerRole.playerId)
		}
	}

	players := findPlayersById(db, playerIds)
	playerByIdMap := make(map[int64]Player)
	for _, player := range players {
		playerByIdMap[player.id] = player
	}

	var matchResults = make(map[int64]structs.Result)
	for _, matchId := range goodMatchIds {
		matchResults[matchId] = makeResult(game, playerByIdMap, matchTeamResultsByMatchId[matchId], matchPlayerRolesByMatchId[matchId])
	}
	return matchResults
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
