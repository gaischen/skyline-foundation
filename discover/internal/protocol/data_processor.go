package protocol

//处理服务端数据的定义
//goland:noinspection ALL
type ServerDataProcessor interface {
	Online(meta ServiceMeta) error
	Offline(meta ServiceMeta) error
	Register(meta ServiceMeta) error
	Remove(meta ServiceMeta) error
}

type ClientDataProcessor interface {

}
