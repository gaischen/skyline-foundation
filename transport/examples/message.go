package examples

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	length uint32
	id     uint32
	flag   uint32
	value  string
}

func Msg2Bytes(m *Message) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(0))
	binary.Write(buf, binary.BigEndian, m.id)
	binary.Write(buf, binary.BigEndian, m.flag)
	buf.WriteString(m.value)
	bytes := buf.Bytes()
	length := len(buf.Bytes())
	binary.BigEndian.PutUint32(bytes[0:4], uint32(length))
	return bytes
}

func Bytes2Msg(data []byte) *Message {
	msg := &Message{}
	msg.length = binary.BigEndian.Uint32(data[0:4])
	msg.id = binary.BigEndian.Uint32(data[4:8])
	msg.flag = binary.BigEndian.Uint32(data[8:12])
	le := len(data)
	msg.value = string(data[12:le])
	return msg
}
