package tg

import (
	"database/sql"
	"log"

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
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "register" {
			player := structs.Player{Name: update.Message.From.UserName, TgId: update.Message.From.ID}
			RegisterPlayer(bot, update.Message.Chat.ID, update.Message.MessageID, sqlDb, player)
		}
	}
}
