package nap

import "errors"

type FrameType uint8
const (
	FrameDate FrameType = 0x0
	FramePing FrameType = 0x1
)

type Flags uint8
const (
	FlagPingAck Flags = 0x1
)

func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

var MaxFrameSize = uint32(64*1024)

var ErrFrameTooLarge = errors.New("frame too large")

var (
	FrameHeaderSize = 10
)
type FrameHeader struct {
	Type FrameType
	Flags Flags
	Length uint32
	ID uint32
}

type Frame struct {
	Header  FrameHeader
	Payload []byte

	err error
}

func (f *Frame) IsPingAck() bool {
	if f.Header.Type == FramePing {
		if f.Header.Flags.Has(FlagPingAck){
			return true
		}
	}

	return false
}
