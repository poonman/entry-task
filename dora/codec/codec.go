package codec

import "strings"

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Name() string
}

var registeredCodecs = make(map[string]Codec)

func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	typ := strings.ToLower(codec.Name())
	registeredCodecs[typ] = codec
}

func GetCodec(typ string) Codec {
	return registeredCodecs[typ]
}
