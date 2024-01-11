dev_dependencies:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/cosmtrek/air@latest

migration_create:
	migrate create -ext sql -dir ./migrations -seq $(name)

migration_up:
	go run ./cmd/migration

migration_drop:
	migrate -database postgres://postgres:price@localhost:5432/postgres?sslmode=disable -path ./migration drop

templ:
	@templ generate

build_sys_arch:
	mkdir -p ../dist/$(system)_$(arch)/catalog;
	cp config.yaml ../dist/$(system)_$(arch)/catalog;
	cp docker-compose.yml ../dist/$(system)_$(arch)/catalog;
	cp -r migrations/ ../dist/$(system)_$(arch)/catalog/migrations;
	GOOS=$(system) GOARCH=$(arch) go build -o ../dist/$(system)_$(arch)/catalog/catalog$(extension) ./cmd/app;
	GOOS=$(system) GOARCH=$(arch) go build -o ../dist/$(system)_$(arch)/catalog/migrate$(extension) ./cmd/migration;


build:
	@if test -z "$(system)" || test -z "$(arch)"; then\
		make build_sys_arch system=windows arch=amd64 extension=.exe;\
		make build_sys_arch system=windows arch=386 extension=.exe;\
		make build_sys_arch system=darwin arch=amd64;\
		make build_sys_arch system=darwin arch=arm64;\
		make build_sys_arch system=linux arch=386;\
		make build_sys_arch system=linux arch=amd64;\
		make build_sys_arch system=linux arch=arm64;\
	else\
		make build_sys_arch system=$(system) arch=$(arch);\
	fi

run:
	IS_DEV=true air
