module example.com/main

go 1.23.4

replace example.com/db => ./db

replace example.com/structs => ./structs

require (
	example.com/db v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.24
)

require example.com/structs v0.0.0-00010101000000-000000000000
