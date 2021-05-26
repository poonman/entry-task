package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/poonman/entry-task/dora/log"
	"io"
	"net"
)

var (
	MaxMessageLength = uint32(64 * 1024) // 64KB
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
			Head:   &head,
			Method: m.PkgHead.Method,
			Meta:   m.PkgHead.Meta,
		},
	}
	return mm
}
func ReadMessage(r io.Reader) (msg *Message, err error) {
	magicBytes := make([]byte, 1)
	pkgLenBytes := make([]byte, 4)

	log.Debugf("magic...")
	_, err = io.ReadFull(r, magicBytes)
	if err != nil {
		return
	}

	log.Debugf("magic:[%v]", magicBytes)

	err = checkMagic(magicBytes[0])
	if err != nil {
		return
	}

	// pkg
	log.Debugf("pkgLen...")
	_, err = io.ReadFull(r, pkgLenBytes)
	if err != nil {
		return
	}

	pkgLen := binary.BigEndian.Uint32(pkgLenBytes)
	log.Debugf("pkgLen:%d", pkgLen)

	// payload
	log.Debugf("payload len...")
	_, err = io.ReadFull(r, pkgLenBytes)
	if err != nil {
		return
	}

	payloadLen := binary.BigEndian.Uint32(pkgLenBytes)
	log.Debugf("payloadLen:%d", payloadLen)

	if pkgLen+payloadLen > MaxMessageLength {
		err = ErrMessageTooLong
		return
	}

	log.Debugf("pkg...")
	pkgBytes := make([]byte, pkgLen)
	_, err = io.ReadFull(r, pkgBytes)
	if err != nil {
		return
	}

	log.Debugf("payload")
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

	log.Debugf("readMessage success...")
	msg = &Message{
		PkgHead: pkg,
		Payload: payloadBytes,
	}

	return
}

func WriteMessage(conn net.Conn, msg *Message) (err error) {

	log.Debugf("WriteMessage begin...")

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

	totalLen := 1 + 4 + 4 + pkgLen + payloadLen
	data := make([]byte, 9, totalLen)
	data[0] = Magic

	binary.BigEndian.PutUint32(data[1:5], pkgLen)
	binary.BigEndian.PutUint32(data[5:9], payloadLen)

	data = append(data, pkgBytes...)
	data = append(data, msg.Payload...)

	log.Debugf("data:[%+v]", data)

	_, err = conn.Write(data)
	//_, err = w.Write(data)
	if err != nil {
		return
	}


	log.Debugf("WriteMessage end...")
	return
}

func checkMagic(magic byte) (err error) {
	if magic != Magic {
		return errors.New("invalid magic")
	}
	return
}
