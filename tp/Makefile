GO_BINARY = go
TARGET_BIN = break-sawtooth-history-tp
TARGET_DIR = ../build

all: build

build:
	${GO_BINARY} build -o ${TARGET_DIR}/${TARGET_BIN} main.go

clean:
	rm -fv ${TARGET_DIR}/${TARGET_BIN}
