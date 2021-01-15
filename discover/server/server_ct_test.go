package server

import (
	"fmt"
	"testing"
)

func TestNewServerCT(t *testing.T) {
	s := NewBasicServer(nil)
	//s.DiscoverType()
	serverCT := NewServerCT(s)
	fmt.Println(serverCT.CurrentDiscoverType)

}
