package protocol

import (
	"bytes"
	"encoding/gob"
	"time"
)

//内部传输协议
type Transport struct {
	length      int
	messageType int // 0:slave->master 1:master->slave 2:client->server
	header      map[string]string
	startTime   time.Time
	endTime     time.Time
	serviceMeta ServiceMeta
}

func Serialize(transport *Transport) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(transport)
	return result.Bytes()
}

func Deserialize(data []byte) *Transport {
	var transport *Transport
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(&transport)
	return transport
}
