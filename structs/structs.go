package structs

type Player struct {
	Name string
	TgId int64
}

type Role struct {
	Name     string
	IsUnique bool
}

type Game struct {
	Name  string
	Roles map[string][]Role // team -> list of roles
}

type Result struct {
	Game        Game
	PlayerRoles map[Player]string
	TeamOrder   []string
}
