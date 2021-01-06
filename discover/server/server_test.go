package server

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewBasicServer(t *testing.T) {
	s := NewBasicServer(nil)
	//s.DiscoverType()
	basicServer, ok := s.(*basicServer)
	if ok {
		fmt.Println(reflect.TypeOf(basicServer))
	}
	fmt.Println(basicServer.DiscoverType())
}

func TestBasicServer_Listen(t *testing.T) {
	s := NewBasicServer(nil)
	//s.DiscoverType()
	basicServer, ok := s.(*basicServer)
	if !ok {
		return
	}
	s, err := basicServer.Listen("localhost", "8080").Start()
	if err != nil {
		fmt.Println(err)
	}


}
