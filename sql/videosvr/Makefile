install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

up:
	goose mysql "root:zl020613@tcp(localhost:3308)/short_video_db?parseTime=true" up -dir .

down:
	goose mysql "root:zl020613@tcp(localhost:3308)/short_video_db?parseTime=true" down -dir .

create:
	goose mysql "root:zl020613@tcp(localhost:3308)/short_video_db?parseTime=true" create ${name} sql -dir .