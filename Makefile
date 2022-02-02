.SILENT:
.PHONY:

run:
	go run cmd/main.go

build:
	go build -o astral cmd/main.go
	go build -o consumer cmd/rabbitmq/consumer.go

migrate-up:
	migrate -path "./migrations" -database "postgres://astral:astral@localhost:5432/astral?sslmode=disable" up

migrate-down:
	migrate -path "./migrations" -database "postgres://astral:astral@localhost:5432/astral?sslmode=disable" down

migrate-create:
	migrate create -ext sql -seq -dir "./migrations" $(filter-out $@,$(MAKECMDGOALS))

.PHONY:
rabbitmq:
	docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

queue:
	go run cmd/rabbitmq/consumer.go