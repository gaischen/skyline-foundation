package network

type networkFeature struct {
	m byte
}

func (n *networkFeature) mask() byte {
	return n.m
}

func (n *networkFeature) isEnable(flag uint8) bool {
	return (flag & n.m) != 0
}

func (n *networkFeature) enable(flag uint8) uint8 {
	return flag | n.m
}

func (n *networkFeature) disable(flag uint8) uint8 {
	return flag & ^n.m
}

var compressFeature = networkFeature{byte(0x80)}
var heartbeatFeature = networkFeature{byte(0x40)}

