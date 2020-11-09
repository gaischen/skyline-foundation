package network

import (
	"sync"
)

type codec interface {
	Read(data []byte, obj interface{}) error
	Write(msg interface{}) ([]byte, error)
}

func init() {
	//codecsMap.setCodec(serializations.Serialization_JSON, &jsonCodec{})
}

var codecsMap = codecs{sync.Map{}}

type codecs struct {
	sync.Map
}

func (c *codecs) getCodec(s string) codec {
	if load, ok := c.Load(s); ok {
		return load.(codec)
	} else {
		return nil
	}
}

func (c *codecs) setCodec(s string, co codec) {
	c.Store(s, co)
}
