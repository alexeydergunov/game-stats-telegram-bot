package structs

import (
	"log"
	"slices"
)

func GetSupportedGames() []Game {
	return []Game{
		{
			Name: "Codenames",
			Roles: map[string][]Role{
				"First": {
					{Name: "FirstTeamCaptain", IsUnique: true},
					{Name: "FirstTeamPlayer", IsUnique: false},
				},
				"Second": {
					{Name: "SecondTeamCaptain", IsUnique: true},
					{Name: "SecondTeamPlayer", IsUnique: false},
				},
			},
		},
		{
			Name: "Decrypto",
			Roles: map[string][]Role{
				"White": {
					{Name: "WhitePlayer", IsUnique: false},
				},
				"Black": {
					{Name: "BlackPlayer", IsUnique: false},
				},
			},
		},
		{
			Name: "Secret Hitler",
			Roles: map[string][]Role{
				"Liberals": {
					{Name: "Liberal", IsUnique: false},
				},
				"Fascists": {
					{Name: "Fascist", IsUnique: false},
					{Name: "Hitler", IsUnique: true},
				},
			},
		},
		{
			Name: "Avalon",
			Roles: map[string][]Role{
				"Knights": {
					{Name: "Merlin", IsUnique: true},
					{Name: "Perceval", IsUnique: true},
					{Name: "Knight", IsUnique: false},
				},
				"Assassins": {
					{Name: "Morgana", IsUnique: true},
					{Name: "Assassin", IsUnique: false},
				},
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
