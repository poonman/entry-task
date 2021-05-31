.PHONY: all clean

OUTPUT=etserver

all: clean
	go build -o ./bin/${OUTPUT} cmd/main.go

idl: clean-idl
	protoc -I/Users/wen.pan/source/go/src/github.com/poonman/entry-task/tiger --go_out=./idl/kv --dora_out=./idl/kv kv.proto

clean-idl:
	rm -rf idl/kv/*

clean:
	rm -f bin/${OUTPUT}