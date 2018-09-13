.PHONY: all example-client attach-server

all: example-client attach-server

example-client:
	if [ ! -d ./vendor ]; then dep ensure -vendor-only; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o _output/example-client ./example-client/main.go
attach-server:
	if [ ! -d ./vendor ]; then dep ensure -vendor-only; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o _output/attach-server ./server/main.go

clean:
	go clean -r -x
	-rm -rf _output
