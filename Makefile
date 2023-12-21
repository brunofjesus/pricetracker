build:
	@cd ./catalog && make build system=$(system) arch=$(arch)
	@cd ./stores/pingodoce  && make build system=$(system) arch=$(arch)
	@cd ./stores/worten  && make build system=$(system) arch=$(arch)