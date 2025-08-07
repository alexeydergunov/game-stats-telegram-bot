package tg

import (
	"database/sql"
	"log"
	"strconv"

	"example.com/structs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunBot(token string, sqlDb *sql.DB, games []structs.Game) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		message := update.Message
		if message == nil {
			continue
		}

		if message.Command() == "register" {
			player := structs.Player{Name: message.From.UserName, TgId: message.From.ID}
			RegisterPlayer(bot, message.Chat.ID, message.MessageID, sqlDb, player)
		}

		if message.Command() == "list_games" {
			ListGames(bot, message.Chat.ID, message.MessageID, games)
		}

		if message.Command() == "list_players" {
			ListPlayers(bot, message.Chat.ID, message.MessageID, sqlDb)
		}

		if message.Command() == "get_match_result" {
			commandArguments := message.CommandArguments()
			matchId, err := strconv.ParseInt(commandArguments, 10, 64)
			if err != nil {
				replyMessage := tgbotapi.NewMessage(message.Chat.ID, "Wrong arguments, use '/get_match_result [match_id]'")
				replyMessage.ReplyToMessageID = message.MessageID
				send(bot, replyMessage)
			} else {
				GetMatchResult(bot, message.Chat.ID, message.MessageID, sqlDb, matchId)
			}
		}
	}
}
