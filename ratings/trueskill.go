package ratings

import (
	"log"
	"slices"

	"example.com/structs"

	"github.com/gami/go-trueskill"
)

func CalcTrueskillRatings(game structs.Game, matchIdResultMap map[int64]structs.Result) map[structs.Player]float64 {
	var matchIds []int64
	for matchId := range matchIdResultMap {
		matchIds = append(matchIds, matchId)
	}
	slices.Sort(matchIds)

	ts := *trueskill.NewTrueSkill()
	var playerRatingMap = make(map[structs.Player]trueskill.Rating)

	for _, matchId := range matchIds {
		matchResult := matchIdResultMap[matchId]
		for player := range matchResult.PlayerRoles {
			_, contains := playerRatingMap[player]
			if !contains {
				playerRatingMap[player] = *ts.CreateRating()
			}
		}
	}

	for _, matchId := range matchIds {
		matchResult := matchIdResultMap[matchId]
		var playerTeams [][]structs.Player
		var ratingGroups [][]*trueskill.Rating
		for _, team := range matchResult.TeamOrder {
			var roles []string
			for _, role := range game.Roles[team] {
				roles = append(roles, role.Name)
			}
			var playerTeam []structs.Player
			var ratingGroup []*trueskill.Rating
			for player, role := range matchResult.PlayerRoles {
				if slices.Contains(roles, role) {
					rating := playerRatingMap[player]
					playerTeam = append(playerTeam, player)
					ratingGroup = append(ratingGroup, &rating)
				}
			}
			playerTeams = append(playerTeams, playerTeam)
			ratingGroups = append(ratingGroups, ratingGroup)
		}
		log.Println("Rating groups before:", ratingGroups)
		newRatingGroups, err := ts.Rate(ratingGroups)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("Rating groups after:", newRatingGroups)
		for i := 0; i < len(playerTeams); i++ {
			for j := 0; j < len(playerTeams[i]); j++ {
				player := playerTeams[i][j]
				playerRatingMap[player] = *newRatingGroups[i][j]
			}
		}
	}

	var playerFloatRatingMap = make(map[structs.Player]float64)
	for player, rating := range playerRatingMap {
		playerFloatRatingMap[player] = ts.Expose(&rating)
	}
	return playerFloatRatingMap
}
