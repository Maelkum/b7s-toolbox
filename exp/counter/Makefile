all: counter wasm

counter:
	go build -o counter
	
wasm:
	GOOS=wasip1 GOARCH=wasm go build -o counter.wasm

clean:
	rm ./counter
	rm ./counter.wasm