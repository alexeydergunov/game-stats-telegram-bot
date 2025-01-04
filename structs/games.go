package structs

func GetSupportedGames() []Game {
	return []Game{
		{
			Name: "Codenames",
			Roles: map[string][]string{
				"First":  {"Captain", "Player"},
				"Second": {"Captain", "Player"},
			},
		},
		{
			Name: "Decrypto",
			Roles: map[string][]string{
				"First":  {"Player"},
				"Second": {"Player"},
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
