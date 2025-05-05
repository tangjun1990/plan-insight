ver=:1.1.1

run:
	go run cmd/server/server.go --config=configs/config_local.toml
build:
	go build -mod=vendor -o build/bin/server cmd/server/server.go
	cp -rf web/index.html build/bin/
	cp -rf configs/config_local.toml build/bin/config.toml
push:
	


# swag文档
swag:
	swag fmt
	swag init -g internal/server/router/swagger.go -o docs