package nap

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

func readFrameHeader(conn net.Conn) (header FrameHeader, err error) {
	buf := make([]byte, FrameHeaderSize)

	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return
	}

	header = FrameHeader{
		Length: binary.BigEndian.Uint32(buf[0:4]),
		ID: binary.BigEndian.Uint32(buf[4:8]),
		Type:   FrameType(buf[8]),
		Flags:  Flags(buf[9]),
	}

	return
}
func  recv(conn net.Conn) (frame *Frame, err error) {
	var header FrameHeader
	var payload []byte

	header, err = readFrameHeader(conn)
	if err != nil {
		return
	}

	if header.Length > MaxFrameSize {
		err = ErrFrameTooLarge
		return
	}

	if header.Length > 0 {
		payload = make([]byte, header.Length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			return
		}
	}

	frame = &Frame{
		Header:  header,
		Payload: payload,
	}

	return
}

func writeFrameHeader(conn net.Conn, header *FrameHeader) (err error) {
	buf := make([]byte, FrameHeaderSize)

	binary.BigEndian.PutUint32(buf[0:4], header.Length)
	binary.BigEndian.PutUint32(buf[4:8], header.ID)
	buf[8] = byte(header.Type)
	buf[9] = byte(header.Flags)

	_, err = conn.Write(buf)
	return
}

func send(conn net.Conn, frame *Frame) (err error) {
	err = writeFrameHeader(conn, &frame.Header)
	if err != nil {
		return
	}

	if frame.Header.Length > 0 {
		_, err = conn.Write(frame.Payload)
		if err != nil {
			return
		}
	}

	return
}


func minTime(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
