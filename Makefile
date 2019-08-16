GO_BINARY = go

.PHONY: tp cli docker

all: tp cli

tp:
	${MAKE} -C tp

cli:
	${MAKE} -C cli

docker:
	docker build -f docker/Dockerfile-Base -t taekion/break-sawtooth-history-base .
	docker build -f docker/Dockerfile-TP -t taekion/break-sawtooth-history-tp .
	docker build -f docker/Dockerfile-CLI -t taekion/break-sawtooth-history-cli .

dep:
	@echo "Installing Go dependencies..."
	go get -d github.com/google/uuid
	go get -d github.com/jessevdk/go-flags
	go get -d github.com/spf13/pflag
	go get -d github.com/taekion-org/sawtooth-client-sdk-go

clean:
	${MAKE} -C tp clean
	${MAKE} -C cli clean
	rm -rfv build/*
