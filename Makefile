.PHONY: bin clean test

.DEFAULT_COAL: bin

bin:
	go build -o build/nsctl main.go

clean:
	rm -rf build

test:
	go build -o build/nsctl main.go
	cp build/nsctl integration-test/nsctl