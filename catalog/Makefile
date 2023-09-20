migration_tool_install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migration_create:
	migrate create -ext sql -dir ./migration -seq $(name)

migration_up:
	migrate -database postgres://postgres:price@localhost:5432/postgres?sslmode=disable -path ./migration up

migration_drop:
	migrate -database postgres://postgres:price@localhost:5432/postgres?sslmode=disable -path ./migration drop