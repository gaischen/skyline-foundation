package protocol

import (
	"fmt"
	"testing"
	"time"
)

func TestSerialize(t *testing.T) {
	transport := &Transport{
		length:      100,
		messageType: 1,
		header:      make(map[string]string),
		startTime:   time.Now(),
		endTime:     time.Now(),
		serviceMeta: ServiceMeta{
			ServiceName: "testServiceName",
		},
	}
	bts := Serialize(transport)
	fmt.Println(bts)

	newTrans := Deserialize(bts)
	fmt.Println(newTrans)
}
