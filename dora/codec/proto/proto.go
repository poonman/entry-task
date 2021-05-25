package proto

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/poonman/entry-task/dora/codec"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

func init() {
	codec.RegisterCodec(baseCodec{})
}

// baseCodec is a Codec implementation with protobuf. It is the default baseCodec
type baseCodec struct{}

func (baseCodec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return proto.Marshal(vv)
}

func (baseCodec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv)
}

func (baseCodec) Name() string {
	return Name
}
