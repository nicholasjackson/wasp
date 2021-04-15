build:
	GOOS=js GOARCH=wasm go build  -o ./module.wasm .

build-tiny-go:
	tinygo build -o ./module.wasm -target wasi -gc leaking .

