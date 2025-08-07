package tg

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"sort"
	"strconv"
	"strings"

	"example.com/db"
	"example.com/structs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func deleteFromSlice[T comparable](arr []T, toDelete []T) []T {
	var result []T
	for _, x := range arr {
		if !slices.Contains(toDelete, x) {
			result = append(result, x)
		}
	}
	return result
}

func RunBot(token string, sqlDb *sql.DB, games []structs.Game) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	commandsConfig := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "register", Description: "Register new player"},
		tgbotapi.BotCommand{Command: "list_games", Description: "List supported games"},
		tgbotapi.BotCommand{Command: "list_players", Description: "List registered players"},
		tgbotapi.BotCommand{Command: "get_match_result", Description: "Get match result (argument [match_id])"},
		tgbotapi.BotCommand{Command: "register_match", Description: "Register new match (with dialog)"},
	)
	_, err = bot.Request(commandsConfig)
	if err != nil {
		log.Panic(err)
	}

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		message := update.Message
		if message == nil {
			continue
		}

		command := message.Command()
		chatId := message.Chat.ID

		if command == "register" {
			player := structs.Player{Name: message.From.UserName, TgId: message.From.ID}
			RegisterPlayer(bot, chatId, message.MessageID, sqlDb, player)
		}

		if command == "list_games" {
			ListGames(bot, chatId, message.MessageID, games)
		}

		if command == "list_players" {
			ListPlayers(bot, chatId, message.MessageID, sqlDb)
		}

		if command == "get_match_result" {
			commandArguments := message.CommandArguments()
			matchId, err := strconv.ParseInt(commandArguments, 10, 64)
			if err != nil {
				replyMessage := tgbotapi.NewMessage(chatId, "Wrong arguments, use '/get_match_result [match_id]'")
				replyMessage.ReplyToMessageID = message.MessageID
				send(bot, replyMessage)
			} else {
				GetMatchResult(bot, chatId, message.MessageID, sqlDb, matchId)
			}
		}

		if command == "register_match" {
			game := enterGame(bot, updates, chatId)
			if game == nil {
				cancelMessage := tgbotapi.NewMessage(chatId, "New match registration was cancelled")
				cancelMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				send(bot, cancelMessage)
				continue
			}
			replyMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("You selected game %s", game.Name))
			replyMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			send(bot, replyMessage)

			allPlayers := db.GetAllPlayers(sqlDb)
			var playerRoles = make(map[structs.Player]string)
			for team, roles := range game.Roles {
				for _, role := range roles {
					players := enterPlayersForRole(bot, updates, chatId, allPlayers, team, role)
					if players == nil {
						cancelMessage := tgbotapi.NewMessage(chatId, "New match registration was cancelled")
						cancelMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						send(bot, cancelMessage)
						continue
					}
					messageText := fmt.Sprintf("You selected the following players for team %s, role %s:\n", team, role)
					for _, player := range *players {
						messageText += fmt.Sprintf("  - %v\n", player)
					}
					for _, player := range *players {
						playerRoles[player] = role
					}
					replyMessage := tgbotapi.NewMessage(chatId, messageText)
					replyMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					send(bot, replyMessage)

					allPlayers = deleteFromSlice(allPlayers, *players)
				}
			}

			var teams []string
			for team := range game.Roles {
				teams = append(teams, team)
			}
			var teamOrder []string
			failure := false
			for place := 1; place <= len(teams); place++ {
				team := enterTeam(bot, updates, chatId, teams, int64(place))
				if team == nil {
					cancelMessage := tgbotapi.NewMessage(chatId, "New match registration was cancelled")
					cancelMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					send(bot, cancelMessage)
					continue
				}
				if len(*team) == 0 {
					cancelMessage := tgbotapi.NewMessage(chatId, "Empty team was selected, cancelling match registration")
					cancelMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					send(bot, cancelMessage)
					continue
				}
				if slices.Contains(teamOrder, *team) {
					replyMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("Team %s was already selected, cancelling match registration", *team))
					replyMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					send(bot, replyMessage)
					failure = true
					break
				}
				teamOrder = append(teamOrder, *team)
			}
			if failure {
				continue
			}

			result := structs.Result{
				Game:        *game,
				PlayerRoles: playerRoles,
				TeamOrder:   teamOrder,
			}
			RegisterMatch(bot, message.Chat.ID, message.MessageID, sqlDb, result)
		}
	}
}

func enterGame(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, chatId int64) *structs.Game {
	supportedGames := structs.GetSupportedGames()
	var gameByNameMap = make(map[string]structs.Game)
	for _, game := range supportedGames {
		gameByNameMap[game.Name] = game
	}
	// TODO multiple rows
	var buttonRows [][]tgbotapi.KeyboardButton
	for game := range gameByNameMap {
		buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(game)))
	}
	buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Cancel")))

	message := tgbotapi.NewMessage(chatId, "Select game:")
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttonRows...)
	bot.Send(message)

	for update := range updates {
		message := update.Message
		if message == nil {
			continue
		}

		text := message.Text
		log.Println("User reply for enterGame():", text)
		if text == "Cancel" {
			log.Println("Cancel received")
			return nil
		}

		game, ok := gameByNameMap[text]
		if !ok {
			log.Println("Couldn't find such game")
			return nil
		}
		return &game
	}
	return nil
}

func enterTeam(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, chatId int64, teams []string, place int64) *string {
	// TODO multiple rows
	var buttonRows [][]tgbotapi.KeyboardButton
	for _, team := range teams {
		buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(team)))
	}
	buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Cancel")))

	message := tgbotapi.NewMessage(chatId, fmt.Sprintf("Select team which took place %d:", place))
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttonRows...)
	bot.Send(message)

	for update := range updates {
		message := update.Message
		if message == nil {
			continue
		}

		text := message.Text
		log.Println("User reply for enterTeam():", text)
		if text == "Cancel" {
			log.Println("Cancel received")
			return nil
		}

		if !slices.Contains(teams, text) {
			log.Println("Couldn't find such game")
			return nil
		}
		return &text
	}
	return nil
}

func enterPlayersForRole(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, chatId int64, allPlayers []structs.Player, team string, role string) *[]structs.Player {
	sort.Slice(allPlayers, func(i int, j int) bool {
		return allPlayers[i].Name < allPlayers[j].Name
	})

	var selectedPlayers []structs.Player

	for {
		// TODO multiple rows
		var buttonRows [][]tgbotapi.KeyboardButton
		for _, player := range allPlayers {
			buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(fmt.Sprintf("%s %d", player.Name, player.TgId))))
		}
		buttonRows = append(buttonRows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Cancel"), tgbotapi.NewKeyboardButton("Finish")))

		message := tgbotapi.NewMessage(chatId, fmt.Sprintf("Select player for team %s, role %s, or press Finish", team, role))
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttonRows...)
		bot.Send(message)

		for update := range updates {
			message := update.Message
			if message == nil {
				continue
			}

			text := message.Text
			log.Println("User reply for enterPlayersForRole():", text)
			if text == "Cancel" {
				log.Println("Cancel received")
				return nil
			}
			if text == "Finish" {
				log.Println("Finish received")
				return &selectedPlayers
			}
			fields := strings.Fields(text)
			if len(fields) != 2 {
				log.Println("Wrong user entered")
				return nil
			}
			name := fields[0]
			tgId, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				log.Fatalln(err.Error())
				return nil
			}
			selectedPlayer := structs.Player{Name: name, TgId: tgId}
			if slices.Contains(selectedPlayers, selectedPlayer) {
				log.Println("Already selected")
				return nil
			}
			if !slices.Contains(allPlayers, selectedPlayer) {
				log.Println("Not avaliable for selection")
				return nil
			}
			selectedPlayers = append(selectedPlayers, selectedPlayer)
			allPlayers = deleteFromSlice(allPlayers, []structs.Player{selectedPlayer})
			messageText := fmt.Sprintf("Player %v was selected for team %s, role %s", selectedPlayer, team, role)
			log.Println(messageText)
			replyMessage := tgbotapi.NewMessage(chatId, messageText)
			replyMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(replyMessage)
			break
		}
	}
}
