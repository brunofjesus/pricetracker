build_sys_arch:
	mkdir -p ../../dist/$(system)_$(arch)/store/pcdiga/
	GOOS=$(system) GOARCH=$(arch) go build -o ../../dist/$(system)_$(arch)/store/pcdiga/pcdiga$(extension) main.go
	cp config.yaml ../../dist/$(system)_$(arch)/store/pcdiga/

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
