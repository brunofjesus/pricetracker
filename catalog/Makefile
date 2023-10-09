migration_tool_install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migration_create:
	migrate create -ext sql -dir ./migration -seq $(name)

migration_up:
	migrate -database postgres://postgres:price@localhost:5432/postgres?sslmode=disable -path ./migration up

migration_drop:
	migrate -database postgres://postgres:price@localhost:5432/postgres?sslmode=disable -path ./migration drop

build:
	mkdir -p ../dist/catalog
	cp config.yaml ../dist/catalog
	cp docker-compose.yml ../dist/catalog
	cd src && go build -o ../../dist/catalog/catalog cmd/main.go