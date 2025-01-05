module example.com/main

go 1.23.4

replace example.com/db => ./db

replace example.com/structs => ./structs

replace example.com/tg => ./tg

require (
	example.com/db v0.0.0-00010101000000-000000000000
	example.com/structs v0.0.0-00010101000000-000000000000
	example.com/tg v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.24
)

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
