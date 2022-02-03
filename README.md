# Astral

This is just learning project. Go http server application.

## Features
- Go
- REST API
- PostgreSQL
- RabbitMQ
- Redis
- Clean Architecture
- Graceful Shutdown
- Migrations
- Mail
- JWT authentication
- Configuration

## Installation
To configure the application, follow these steps:
1. Create and configure **.env** file in the root directory
2. Configure server config in the **configs** folder(Don't forget to configure database)
3. Change database DSN in the **Makefile**
4. Run database migrations:
```bash
$ make migrate-up
```
5. Run RabbitMQ docker container:
```bash
$ make rabbitmq
```
6. Run http server:
```bash
$ make run
```