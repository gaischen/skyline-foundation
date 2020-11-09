package examples

import (
	"reflect"
	"unsafe"
)

type Message struct {
	length int
	id     uint32
	flag   int
	value  string
}

var sizeOfMyStruct = int(unsafe.Sizeof(Message{}))

func MessageToBytes(s *Message) []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfMyStruct
	x.Cap = sizeOfMyStruct
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func BytesToMessage(b []byte) *Message {
	return (*Message)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}
