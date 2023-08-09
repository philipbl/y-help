default: bin/y-help bin/y-help-handler bin/y-help-linux bin/y-help-handler-linux

bin/y-help: cmd/y-help/main.go
	go build -o $@ $^

bin/y-help-linux: cmd/y-help/main.go
	GOOS=linux GOARCH=amd64 go build -o $@ $^

bin/y-help-handler: cmd/y-help-handler/main.go
	go build -o $@ $^

bin/y-help-handler-linux: cmd/y-help-handler/main.go
	GOOS=linux GOARCH=amd64 go build -o $@ $^

clean:
	rm -rf bin/y-help bin/y-help-handler bin/y-help-linux bin/y-help-handler-linux
