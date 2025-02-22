setup: 
	go mod tidy
run:
	go run main.go
clean:
	rm -f application
build: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o application .

docker: build
	docker build -t apigateway .

dockerrun:
    # run with container name
	docker container stop apigateway
	docker container rm apigateway
	docker run -p 8000:8000 -d --name apigateway apigateway