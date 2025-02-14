.PHONY: test
BUILD_VERSION=$(or ${VERSION}, dev)

build: test lint

test:
	go test ./...

lint:
	golangci-lint run

linux_amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/gcsdeploy-linux-amd64 ./cmd/gcsdeploy
	cd ./dist/ && cp ./gcsdeploy-linux-amd64 ./gcsdeploy && tar cfz gcsdeploy-linux-amd64.tar.gz ./gcsdeploy

linux_arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/gcsdeploy-linux-arm64 ./cmd/gcsdeploy
	cd ./dist/ && cp ./gcsdeploy-linux-arm64 ./gcsdeploy && tar cfz gcsdeploy-linux-arm64.tar.gz ./gcsdeploy

darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/gcsdeploy-darwin-amd64 ./cmd/gcsdeploy
	cd ./dist/ && cp ./gcsdeploy-darwin-amd64 ./gcsdeploy && tar cfz gcsdeploy-darwin-amd64.tar.gz ./gcsdeploy

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/gcsdeploy-darwin-arm64 ./cmd/gcsdeploy
	cd ./dist/ && cp ./gcsdeploy-darwin-arm64 ./gcsdeploy && tar cfz gcsdeploy-darwin-arm64.tar.gz ./gcsdeploy

artifacts: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
