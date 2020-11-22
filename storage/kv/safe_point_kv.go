package kv

import "go.etcd.io/etcd/mvcc/mvccpb"

type SafePointKV interface {
	Put(k string, v string) error
	Get(k string) (string, error)
	GetWithPrefix(k string) ([]*mvccpb.KeyValue, error)
}


