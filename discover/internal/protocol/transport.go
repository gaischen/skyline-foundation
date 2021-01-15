package protocol

import "time"

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

	return nil
}

func Deserialize(data []byte) *Transport {

	return nil
}
