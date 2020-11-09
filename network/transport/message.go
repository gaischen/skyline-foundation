package transport

import (
	"bytes"
	"encoding/binary"
	"sync"
)



var (
	messagePool sync.Pool
)

func acquireMessageObject() *Message {
	res := messagePool.Get()
	if res == nil {
		return new(Message)
	}
	return res.(*Message)
}

/**
```
  message protocol
 +-----------------------------------------------+
 |                 Length (32)                   |
 +---------------+---------------+---------------+
 |Version(8)|protocol(8)|    sequence(16)        |
 +-+-------------+---------------+---------------+
 | Error(8) |  flag(8)  |
 +-+-------------+---------------+----------------
 |                 PackageId(32)                 |
 +---------------+---------------+----------------
 |                 ContentLength(32)             |
 +================================================
 |                      						 |
 |					 Content            		 |
 |  											 |
 |================================================
```
*/

type Message struct {
	Length        uint32
	Version       uint8
	Protocol      uint8
	Sequence      uint16
	Error         uint8
	Flag          uint8
	ContentLength uint32
	PackageId     uint32

	Header map[string]string

	ErrorMessage string
	//content is decode by Body byte array through serializations interface
	Content interface{}
	Body    []byte
}

//parse  msg from byte array,not decode body data,
//decode body data use codec interface
func decodeTeslaMsg(buf *ByteBuf, length uint32) (*Message, error) {
	data := buf.Bytes()
	msg := acquireMessageObject() //todo 这边会依赖到connection里面的方法
	msg.Length = length
	offset := 0
	msg.Version = data[offset]
	offset++
	msg.Protocol = data[offset]
	offset++
	msg.Sequence = binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2
	msg.Error = data[offset]
	offset++
	msg.Flag = data[offset]
	offset++
	msg.PackageId = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	msg.ContentLength = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	msg.Body = data[offset:]
	releaseRecvBuf(buf)
	return msg, nil
}

//encode  msg to byte array,you should encode body data before call this method!
//encode body data use codec interface
func encodeMsg(msg *Message) (*bytes.Buffer, error) {

	var (
		buf    *bytes.Buffer
		offset int32
		err    error
	)

	buf = acquireWriteBuf()
	offset = 0

	//length
	//此处只做占位 待数据写入完毕后修改此处的值
	err = binary.Write(buf, binary.BigEndian, uint32(0))
	if err != nil {
		return nil, err
	}
	offset += 4
	//version
	buf.WriteByte(msg.Version)
	offset++
	buf.WriteByte(msg.Protocol)
	offset++
	err = binary.Write(buf, binary.BigEndian, msg.Sequence)
	if err != nil {
		return nil, err
	}
	offset += 2
	buf.WriteByte(msg.Error)
	offset++
	buf.WriteByte(msg.Flag)
	offset++
	err = binary.Write(buf, binary.BigEndian, msg.PackageId)
	if err != nil {
		return nil, err
	}
	offset += 4
	err = binary.Write(buf, binary.BigEndian, uint32(len(msg.Body)))
	if err != nil {
		return nil, err
	}
	offset += 4
	buf.Write(msg.Body)

	bytes := buf.Bytes()
	//把length给补上
	length := len(buf.Bytes())
	binary.BigEndian.PutUint32(bytes[0:4], uint32(length))
	msg.Length = uint32(length)
	return buf, err
}
