package codec

import "strings"

type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation.
	Name() string
}

var registeredCodecs = make(map[string]Codec)

// RegisterCodec registers the provided Codec
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

// GetCodec gets a registered Codec by serialize-type
func GetCodec(typ string) Codec {
	return registeredCodecs[typ]
}
