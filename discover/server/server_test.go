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
	bsrv, ok := s.(*basicServer)
	if !ok {
		return
	}
	bsrv = bsrv.Listen("tcp", ":8080").Start().(*basicServer)
	fmt.Println("server start....")
	bsrv.wg.Wait()
}
