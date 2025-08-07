package structs

import (
	"log"
	"slices"
)

func GetSupportedGames() []Game {
	return []Game{
		{
			Name: "Codenames",
			Roles: map[string][]string{
				"First":  {"FirstTeamCaptain", "FirstTeamPlayer"},
				"Second": {"SecondTeamCaptain", "SecondTeamPlayer"},
			},
		},
		{
			Name: "Decrypto",
			Roles: map[string][]string{
				"White": {"WhitePlayer"},
				"Black": {"BlackPlayer"},
			},
		},
		{
			Name: "Secret Hitler",
			Roles: map[string][]string{
				"Liberals": {"Liberal"},
				"Fascists": {"Fascist", "Hitler"},
			},
		},
		{
			Name: "Avalon",
			Roles: map[string][]string{
				"Knights":   {"Merlin", "Perceval", "Knight"},
				"Assassins": {"Morgana", "Assassin"},
			},
		},
	}
}

func FindGameByName(name string) *Game {
	games := GetSupportedGames()
	for _, game := range games {
		if game.Name == name {
			log.Println("Found game with name", name)
			return &game
		}
	}
	log.Fatalln("Game with name", name, "is not supported")
	return nil
}

func GetTeamByRole(teamRolesMap map[string][]string, roleToFind string) *string {
	for team, roles := range teamRolesMap {
		if slices.Contains(roles, roleToFind) {
			return &team
		}
	}
	log.Fatalln("Couldn't find team for role", roleToFind)
	return nil
}
