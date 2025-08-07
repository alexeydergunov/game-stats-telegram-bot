package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"example.com/db"
	"example.com/structs"
	"example.com/tg"

	_ "github.com/mattn/go-sqlite3"
)

const DB_FILENAME = "./_sqlite3db.bin"

func main() {
	token := os.Getenv("TG_TOKEN")
	if len(token) == 0 {
		tokenBinary, err := os.ReadFile("token.txt")
		if err != nil {
			log.Fatalln(err.Error())
		}
		token = string(tokenBinary)
	}
	if len(token) == 0 {
		log.Fatalln("Token is not set")
	}
	log.Println("Token:", token[:10]+"..."+token[len(token)-10:])

	games := structs.GetSupportedGames()
	for _, game := range games {
		log.Println("Game:", game)
	}

	if fileInfo, err := os.Stat(DB_FILENAME); errors.Is(err, os.ErrNotExist) {
		log.Println("Db file", DB_FILENAME, "does not exist, creating...")
		file, err := os.Create(DB_FILENAME)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("Created db file", DB_FILENAME)
	} else {
		log.Println("Db file", DB_FILENAME, "already exists, size =", fileInfo.Size())
	}
	sqlDb, _ := sql.Open("sqlite3", DB_FILENAME)
	defer sqlDb.Close()

	db.CreateTables(sqlDb)

	// TODO delete
	// fakePlayers := []structs.Player{
	// 	{Name: "fake_1", TgId: 1001},
	// 	{Name: "fake_2", TgId: 1002},
	// 	{Name: "fake_3", TgId: 1003},
	// 	{Name: "fake_4", TgId: 1004},
	// 	{Name: "fake_5", TgId: 1005},
	// }
	// for _, fakePlayer := range fakePlayers {
	// 	db.GetOrInsertPlayer(sqlDb, fakePlayer)
	// }
	// db.InsertMatchResult(sqlDb, structs.Result{
	// 	Game: *structs.FindGameByName("Secret Hitler"),
	// 	PlayerRoles: map[structs.Player]string{
	// 		fakePlayers[0]: "Liberal",
	// 		fakePlayers[1]: "Fascist",
	// 		fakePlayers[2]: "Liberal",
	// 		fakePlayers[3]: "Liberal",
	// 		fakePlayers[4]: "Hitler",
	// 	},
	// 	TeamOrder: []string{"Liberals", "Fascists"},
	// })

	tg.RunBot(token, sqlDb, games)
}
