
all:
	go get github.com/satori/go.uuid	
	GOOS=js GOARCH=wasm go build -o docs/main.wasm
