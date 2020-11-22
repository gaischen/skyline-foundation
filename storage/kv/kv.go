package kv

type kvStore struct {
	clusterId uint64
	uuid      string
	client    Client
}
