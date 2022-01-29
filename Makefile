.SILENT:
.PHONY:

run:
	go run cmd/main.go

build:
	go build -o astral cmd/main.go

migrate-up:
	migrate -path "./migrations" -database "postgres://astral:astral@localhost:5432/astral?sslmode=disable" up

migrate-down:
	migrate -path "./migrations" -database "postgres://astral:astral@localhost:5432/astral?sslmode=disable" down

migrate-create:
	migrate create -ext sql -seq -dir "./migrations" $(filter-out $@,$(MAKECMDGOALS))