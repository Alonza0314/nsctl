.PHONY: bin clean

.DEFAULT_COAL: bin

bin:
	go build -o build/nsctl main.go

clean:
	rm -rf build