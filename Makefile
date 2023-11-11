.PHONY: mocks docker-image test svr setup

all : mocks


mocks:
	mockgen -destination=mocks/mock_das.go -package=mocks -source=das/das.go

docker-image:
	docker build -t dino_svr:latest -f Dockerfile .

test:
	go test -v

setup:
	go get github.com/gorilla/mux
	go get github.com/lib/pq
	go mod tidy

svr:
	go build -o svr main.go handlers.go species.go gen_map.go

lint:
	golangci-lint run *.go

