package structs

type Player struct {
	Name string
	TgId int64
}

type Game struct {
	Name  string
	Roles map[string][]string // team -> list of roles
}

type Result struct {
	Game        Game
	PlayerRoles map[Player]string
	TeamOrder   []string
}
