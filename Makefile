GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

default: clean linux-build

linux-build:
	 			CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/ix-ad-service -v

clean:
	@echo "+ $@"
	rm -rf build
	mkdir -p build

 test:
      go test -v ./...