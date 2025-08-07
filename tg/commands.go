package tg

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"sort"

	"example.com/db"
	"example.com/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func send(bot *tgbotapi.BotAPI, message tgbotapi.MessageConfig) {
	_, err := bot.Send(message)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func RegisterPlayer(bot *tgbotapi.BotAPI, chatId int64, requestMessageId int, sqlDb *sql.DB, player structs.Player) {
	dbPlayerNullable := db.FindOnePlayerByTgId(sqlDb, player.TgId)
	var dbPlayer db.Player
	var messageText string
	if dbPlayerNullable == nil {
		log.Println("Player with tgId", player.TgId, "not found, will register")
		dbPlayer = db.GetOrInsertPlayer(sqlDb, player)
		messageText = fmt.Sprintf("Registered player %v", dbPlayer)
	} else {
		dbPlayer = *dbPlayerNullable
		messageText = fmt.Sprintf("Player %v is already registered", dbPlayer)
	}
	log.Println(messageText)

	message := tgbotapi.NewMessage(chatId, messageText)
	message.ReplyToMessageID = requestMessageId
	send(bot, message)
}

func ListGames(bot *tgbotapi.BotAPI, chatId int64, requestMessageId int, games []structs.Game) {
	messageText := "Supported games:\n"
	for _, game := range games {
		messageText += fmt.Sprintf("- %s\n", game.Name)
		for team, roles := range game.Roles {
			messageText += fmt.Sprintf("  - %s : %v\n", team, roles)
		}
	}
	message := tgbotapi.NewMessage(chatId, messageText)
	message.ReplyToMessageID = requestMessageId
	send(bot, message)
}

func ListPlayers(bot *tgbotapi.BotAPI, chatId int64, requestMessageId int, sqlDb *sql.DB) {
	players := db.FindAllPlayers(sqlDb)
	messageText := "Registered players:\n"
	for _, player := range players {
		messageText += fmt.Sprintf("  - %v\n", player)
	}
	message := tgbotapi.NewMessage(chatId, messageText)
	message.ReplyToMessageID = requestMessageId
	send(bot, message)
}

func GetMatchResult(bot *tgbotapi.BotAPI, chatId int64, requestMessageId int, sqlDb *sql.DB, matchId int64) {
	matchResultNullable := db.GetMatchResult(sqlDb, matchId)
	var messageText string
	if matchResultNullable == nil {
		messageText += fmt.Sprintf("Couldn't find match %d in db", matchId)
	} else {
		matchResult := *matchResultNullable
		messageText += fmt.Sprintf("Found match %d\n", matchId)
		messageText += fmt.Sprintf("  - game: %s\n", matchResult.Game.Name)
		messageText += fmt.Sprintf("  - team order: %v\n", matchResult.TeamOrder)
		var players []structs.Player
		for player := range matchResult.PlayerRoles {
			players = append(players, player)
		}
		game := structs.FindGameByName(matchResult.Game.Name)
		var roleTeamMap = make(map[string]string)
		for team, roles := range game.Roles {
			for _, role := range roles {
				roleTeamMap[role] = team
			}
		}
		sort.Slice(players, func(i int, j int) bool {
			role1 := matchResult.PlayerRoles[players[i]]
			role2 := matchResult.PlayerRoles[players[j]]
			team1 := roleTeamMap[role1]
			team2 := roleTeamMap[role2]
			if team1 != team2 {
				index1 := slices.Index(matchResult.TeamOrder, team1)
				index2 := slices.Index(matchResult.TeamOrder, team2)
				log.Println("Comparing...", role1, role2, team1, team2, index1, index2)
				return index1 < index2
			}
			teamRoles := game.Roles[team1]
			index1 := slices.Index(teamRoles, role1)
			index2 := slices.Index(teamRoles, role2)
			if index1 != index2 {
				log.Println("Comparing...", role1, role2, team1, team2, index1, index2)
				return index1 < index2
			}
			return players[i].Name < players[j].Name
		})
		for _, player := range players {
			messageText += fmt.Sprintf("    - player %s has role %s\n", player.Name, matchResult.PlayerRoles[player])
		}
	}

	message := tgbotapi.NewMessage(chatId, messageText)
	message.ReplyToMessageID = requestMessageId
	send(bot, message)
}

func RegisterMatch(bot *tgbotapi.BotAPI, chatId int64, requestMessageId int, sqlDb *sql.DB, matchResult structs.Result) {
	var notFoundPlayers []structs.Player
	for player := range matchResult.PlayerRoles {
		dbPlayerNullable := db.FindOnePlayerByTgId(sqlDb, player.TgId)
		if dbPlayerNullable == nil {
			log.Println("Player with tgId", player.TgId, "not found, they must register first")
			notFoundPlayers = append(notFoundPlayers, player)
		}
	}

	var messageText string
	if len(notFoundPlayers) > 0 {
		messageText += fmt.Sprintf("Found %d not registered players:\n", len(notFoundPlayers))
		for _, notFoundPlayer := range notFoundPlayers {
			messageText += fmt.Sprintf("  - %v\n", notFoundPlayer)
		}
	} else {
		matchId := db.InsertMatchResult(sqlDb, matchResult)
		messageText += fmt.Sprintf("Registered match %d", matchId)
	}

	message := tgbotapi.NewMessage(chatId, messageText)
	message.ReplyToMessageID = requestMessageId
	send(bot, message)
}
