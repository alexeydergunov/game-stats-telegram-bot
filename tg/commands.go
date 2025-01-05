package tg

import (
	"database/sql"
	"fmt"
	"log"

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
