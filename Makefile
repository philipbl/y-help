default: bin/y-help bin/y-help-server bin/y-help-linux bin/y-help-server-linux

bin/y-help: cmd/y-help/main.go
	go build -o $@ $^

bin/y-help-linux: cmd/y-help/main.go
	GOOS=linux GOARCH=amd64 go build -o $@ $^

bin/y-help-server: cmd/y-help-server/main.go
	go build -o $@ $^

bin/y-help-server-linux: cmd/y-help-server/main.go
	GOOS=linux GOARCH=amd64 go build -o $@ $^

clean:
	rm -rf bin/y-help bin/y-help-server bin/y-help-linux bin/y-help-server-linux
