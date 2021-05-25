package protocol


import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
)

var (
	MaxMessageLength = uint32(64*1024) // 64KB
)

var (
	ErrMessageTooLong = errors.New("message is too long")
)

var (
	Magic byte = 'G'
)


type Message struct {
	PkgHead *PkgHead
	Payload []byte
}

func (m *Message) SetMessageType(typ Head_MessageType) {
	m.PkgHead.Head.MessageType = typ
}

func (m *Message) Clone() *Message {
	head := *m.PkgHead.Head

	mm := &Message{
		PkgHead: &PkgHead{
			Head:                 &head,
			Method:               m.PkgHead.Method,
			Meta:                 m.PkgHead.Meta,
		},
	}
	return mm
}
func ReadMessage(r io.Reader) (msg *Message, err error) {
	magicBytes := make([]byte, 1)
	lenBytes := make([]byte, 4)

	_, err = io.ReadFull(r, magicBytes)
	if err != nil {
		return
	}

	err = checkMagic(magicBytes[0])
	if err != nil {
		return
	}

	// pkg
	_, err = io.ReadFull(r, lenBytes)
	if err != nil {
		return
	}

	pkgLen := binary.BigEndian.Uint32(lenBytes)

	// payload
	_, err = io.ReadFull(r, lenBytes)
	if err != nil {
		return
	}

	payloadLen := binary.BigEndian.Uint32(lenBytes)

	if pkgLen+payloadLen > MaxMessageLength {
		err = ErrMessageTooLong
		return
	}

	pkgBytes := make([]byte, pkgLen)
	_, err = io.ReadFull(r, pkgBytes)
	if err != nil {
		return
	}

	payloadBytes := make([]byte, payloadLen)
	_, err = io.ReadFull(r, payloadBytes)
	if err != nil {
		return
	}

	pkg := &PkgHead{}

	err = proto.Unmarshal(pkgBytes, pkg)
	if err != nil {
		return
	}

	msg = &Message{
		PkgHead: pkg,
		Payload: payloadBytes,
	}

	return
}

func WriteMessage(w io.Writer, msg *Message) (err error) {

	pkgBytes, err := proto.Marshal(msg.PkgHead)
	if err != nil {
		return
	}

	pkgLen := uint32(len(pkgBytes))
	payloadLen := uint32(len(msg.Payload))

	if pkgLen+payloadLen > MaxMessageLength {
		err = ErrMessageTooLong
		return
	}

	totalLen := 1+4+4+pkgLen+payloadLen
	data := make([]byte, 9, totalLen)
	data[0] = Magic

	binary.BigEndian.PutUint32(data[1:5], pkgLen)
	binary.BigEndian.PutUint32(data[5:9], payloadLen)

	data = append(data, pkgBytes...)
	data = append(data, msg.Payload...)

	_, err = w.Write(data)
	return
}


func checkMagic(magic byte) (err error) {
	if magic != Magic {
		return errors.New("invalid magic")
	}
	return
}
