package db

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {

	{
		statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS player (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"name" TEXT,
			"tg_id" INTEGER
		)`)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
		log.Println("player table created")
	}

	{
		statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS match (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"game" TEXT NOT NULL
		)`)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
		log.Println("match table created")
	}

	{
		statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS match_player_role (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"match_id" INTEGER NOT NULL,
			"player_id" INTEGER NOT NULL,
			"role" TEXT NOT NULL
		)`)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
		log.Println("match_player_role table created")
	}

	{
		statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS match_team_result (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"match_id" INTEGER NOT NULL,
			"team" TEXT NOT NULL,
			"place" TEXT NOT NULL
		)`)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
		log.Println("match_team_result table created")
	}
}
