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
	mu        sync.RWMutex      // the gatekeeper for all fields
	file      *os.File          // the underlying file
	buf       []byte            // a buffer to write to
	keys      *btree.BTree      // a tree of all item ordered by key
	exps      *btree.BTree      // a tree of items ordered by expiration
	idxs      map[string]*index // the index trees.
	exmgr     bool              // indicates that expires manager is running.
	flushes   int               // a count of the number of disk flushes
	closed    bool              // set when the database has been closed
	config    Config            // the database configuration
	persist   bool              // do we write to disk
	shrinking bool              // when an aof shrink is in-process.
	lastaofsz int               // the size of the last shrink aof size
}