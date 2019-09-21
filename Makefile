
all:
	go get github.com/google/uuid
	GOOS=js GOARCH=wasm go build -o docs/main.wasm
