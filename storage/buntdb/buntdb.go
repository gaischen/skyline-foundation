package buntdb

type SyncPolicy int

const (
	SYNC_ALLA SyncPolicy = iota
)

type Config struct {
	SyncPolicy SyncPolicy
	AutoShrinkPercentage int
}

type DB struct {

}