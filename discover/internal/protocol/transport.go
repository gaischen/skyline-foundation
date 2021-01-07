package protocol

import "time"

//内部传输协议
type Transport struct {
	length      int
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
