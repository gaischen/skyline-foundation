package network



type RpcSpace string

type JavaObject interface {
	SetRpcSpace()
}

type Invocation struct {
	RpcSpace

	MethodName string

	Arguments []interface{}

	Service string

	Version string

	Group string

	UniqueServiceName string

	Attachments map[string]interface{}
}

func NewInvocation(methodName string, arguments []interface{}, service string, version string, group string, uniqueServiceName string) *Invocation {
	return &Invocation{MethodName: methodName,
		Arguments:         arguments,
		Service:           service,
		Version:           version,
		Group:             group,
		UniqueServiceName: uniqueServiceName,
		Attachments:       make(map[string]interface{}),
	}
}

func (i *Invocation) AddAttachment(key string, v interface{}) {
	if i.Attachments == nil {
		i.Attachments = make(map[string]interface{})
	}
	i.Attachments[key] = v
}

func (i *Invocation) SetTeslaSpace() {
	i.RpcSpace = "com.rpc.core.Invocation"
}

type Result struct {
	Status int32

	Value interface{}

	Desc string

	Exception interface{}

	Attachments map[string]interface{}
}
