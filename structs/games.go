package structs

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
