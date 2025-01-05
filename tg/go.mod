module example.com/tg

go 1.23.4

replace example.com/db => ../db

replace example.com/structs => ../structs

require (
	example.com/db v0.0.0-00010101000000-000000000000
	example.com/structs v0.0.0-00010101000000-000000000000
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
)
