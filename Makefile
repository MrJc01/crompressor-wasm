.PHONY: build serve clean

build:
	GOOS=js GOARCH=wasm go build -o www/crompressor.wasm ./cmd/wasm
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" www/

serve: build
	cd www && python3 -m http.server 8080

clean:
	rm -f www/crompressor.wasm www/wasm_exec.js
