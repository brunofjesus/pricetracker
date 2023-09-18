migration_tool_install:
	go install -tags 'postgres,sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migration_create:
	migrate create -ext sql -dir ./migration/sql -seq $(name)