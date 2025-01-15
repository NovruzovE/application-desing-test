run:
	go run -ldflags "-X main.version=v0.0.1" cmd/main.go

test:
	go test -v ./...

fmt:
	go fmt ./...
