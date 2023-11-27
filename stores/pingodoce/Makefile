build_sys_arch:
	mkdir -p ../../dist/$(system)_$(arch)/store/pingodoce/
	GOOS=$(system) GOARCH=$(arch) go build -o ../../dist/$(system)_$(arch)/store/pingodoce/pingodoce$(extension) main.go
	cp config.yaml ../../dist/$(system)_$(arch)/store/pingodoce/

build:
	make build_sys_arch system=windows arch=amd64 extension=.exe
	make build_sys_arch system=windows arch=386 extension=.exe
	make build_sys_arch system=darwin arch=amd64
	make build_sys_arch system=darwin arch=arm64
	make build_sys_arch system=linux arch=386
	make build_sys_arch system=linux arch=amd64
	make build_sys_arch system=linux arch=arm64
