module example.com/main

go 1.23.4

replace example.com/db => ./db

replace example.com/ratings => ./ratings

replace example.com/structs => ./structs

replace example.com/tg => ./tg

require (
	example.com/db v0.0.0-00010101000000-000000000000
	example.com/structs v0.0.0-00010101000000-000000000000
	example.com/tg v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.24
)

require (
	example.com/ratings v0.0.0-00010101000000-000000000000 // indirect
	github.com/chobie/go-gaussian v0.0.0-20150107165016-53c09d90eeaf // indirect
	github.com/gami/go-trueskill v0.0.0-20210522135225-3b7dd85a7d62 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
)
