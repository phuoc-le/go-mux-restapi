* Start app
Copy .env-example into .env
``go mod download``
``go build -o out/go-mux-restapi``
``./out/go-mux-restapi``


migrate create -ext sql -dir db/migrations -seq create_users_table

migrate -source file://db/migrations/ -database postgres://admin:123456@localhost:5432/test?sslmode=disable up 2

