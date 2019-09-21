
all:
	GOOS=js GOARCH=wasm go build -o docs/main.wasm
