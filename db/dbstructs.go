package db

type Player struct {
	id   int64
	name string
	tgId int64
}

type Match struct {
	id   int64
	game string
}

type MatchPlayerRole struct {
	id       int64
	matchId  int64
	playerId int64
	role     string
}

type MatchTeamResult struct {
	id      int64
	matchId int64
	team    string
	place   int64
}
