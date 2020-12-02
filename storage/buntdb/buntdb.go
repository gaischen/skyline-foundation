package buntdb

type SyncPolicy int

const (
	SYNC_ALLA SyncPolicy = iota
)

type Config struct {
	SyncPolicy SyncPolicy

}

type DB struct {

}