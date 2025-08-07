module example.com/ratings

go 1.23.4

replace example.com/structs => ../structs

require (
	example.com/structs v0.0.0-00010101000000-000000000000
	github.com/gami/go-trueskill v0.0.0-20210522135225-3b7dd85a7d62
)

require github.com/chobie/go-gaussian v0.0.0-20150107165016-53c09d90eeaf // indirect
