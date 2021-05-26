.PHONY: all clean idl clean-idl

OUTPUT=client

all: clean
	go build -o ./bin/${OUTPUT} main.go

idl: clean-idl
	protoc -I/Users/wen.pan/source/go/src/github.com/poonman/entry-task/tiger --go_out=./idl/kv --dora_out=./idl/kv kv.proto

clean:
	rm -f bin/${OUTPUT}

clean-idl:
	rm -rf idl/kv/*