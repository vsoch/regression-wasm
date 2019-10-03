
all:
	go get github.com/sajari/regression
	GOOS=js GOARCH=wasm go build -o docs/main.wasm
