build: clean
	go build -o application .
build-linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o application .
run:
	go run main.go
test:
	go test -v ./...
clean:
	rm -f application

docker: build-linux
	docker build -t opengate .
	rm -f application

docker-run: docker
	@echo "Tagging image as latest"
	docker tag gobaserservice opengate:latest
	@echo "removing existing container named opengate if any"
	docker rm -f opengate || true
	@echo "Running image with name opengate, mapping ports 8085:8085 and 8086:8086"
	docker run -d --name gobaserservice -p 8085:8085 -p 8086:8086 gobaserservice:latest

install: 
# 	go mod tidy
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/gofreego/goutils/cmd/sql-migrator@latest

setup:
	sh ./api/protoc.sh
	go mod tidy

redeploy:
	@echo "Redeploying the application"
	@echo "Pulling latest changes from git"
	git pull
	@echo "Building the docker imamge"
	docker compose build
	@echo "Stopping existing docker containers"
	docker compose down
	@echo "Starting the docker containers"
	docker compose up -d
