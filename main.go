package main

import (
	"database/sql"
	"log"
	"os"

	"example.com/db"
	"example.com/structs"

	_ "github.com/mattn/go-sqlite3"
)

const DB_FILENAME = "./_sqlite3db.bin"

func main() {
	games := structs.GetSupportedGames()
	for _, game := range games {
		log.Println("Game:", game)
	}

	os.Remove(DB_FILENAME)
	log.Println("Creating db file", DB_FILENAME)
	file, err := os.Create(DB_FILENAME)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Created db in file", DB_FILENAME)

	sqlDb, _ := sql.Open("sqlite3", DB_FILENAME)
	defer sqlDb.Close()

	db.CreateTables(sqlDb)

	players := []structs.Player{
		{Name: "Alex", TgId: 11},
		{Name: "Bob", TgId: 22},
		{Name: "Charlie", TgId: 33},
		{Name: "Dennis", TgId: 44},
		{Name: "Eugene", TgId: 55},
	}

	for _, player := range players {
		db.GetOrInsertPlayer(sqlDb, player)
	}

	matchId := db.InsertMatchResult(sqlDb, structs.Result{
		Game: games[3],
		PlayerRoles: map[structs.Player]string{
			players[0]: "Merlin",
			players[1]: "Perceval",
			players[2]: "Knight",
			players[3]: "Morgana",
			players[4]: "Assassin",
		},
		TeamOrder: []string{"Knights", "Assassins"},
	})
	log.Println("Match id:", matchId)

	for _, obj := range db.FindAllPlayers(sqlDb) {
		log.Println("Player:", obj)
	}
	for _, obj := range db.FindAllMatches(sqlDb) {
		log.Println("Match:", obj)
	}
	for _, obj := range db.FindAllMatchPlayerRoles(sqlDb) {
		log.Println("Match player role:", obj)
	}
	for _, obj := range db.FindAllMatchTeamResults(sqlDb) {
		log.Println("Match team result:", obj)
	}
}
