default: bin/y-help bin/y-help-server

bin/y-help: cmd/y-help/main.go
	go build -o $@ $^

bin/y-help-server: cmd/y-help-server/main.go
	go build -o $@ $^

clean:
	rm -rf bin/y-help bin/y-help-server
