.PHONY: bin clean test install uninstall

.DEFAULT_COAL: bin

bin:
	go build -o build/nsctl main.go

clean:
	rm -rf build

test:
	go build -o build/nsctl main.go
	cp build/nsctl integration-test/nsctl

install:
	mv build/nsctl /usr/bin/nsctl

uninstall:
	rm /usr/bin/nsctl