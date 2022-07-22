package kvbench

import (
	"database/sql"
	"os"
	"path"

	"github.com/boltdb/bolt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/flower-corp/rosedb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/peterbourgon/diskv/v3"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tidwall/buntdb"
	"github.com/xujiajun/nutsdb"
	"go.etcd.io/bbolt"
)

// https://github.com/objectbox/objectbox-go
// - transactions
// - column indexing
// - relations
// - queries with filters (lte/gte/prefix/suffix/contains/eq/between/in/notin)
// - aggregations
// - pagination
// - ordering
// - distributed syncing
func NewObjBox(dir string) *objectbox.ObjectBox {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)
	objectBox, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Directory(dir).Build()
	if err != nil {
		panic(err)
	}

	return objectBox
}

// https://github.com/dgraph-io/badger
// - inmemory mode
// - encryption
// - transactions
// - ttl
// - items keys metadata
// - iterating over keys
// - Key-only iteration
// - Prefix scans
func NewBadgerDB(dir string) *badger.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	// opts.ValueLogFileSize = 1024 * 1024

	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/boltdb/bolt
// - nested buckets
// - transactions
// @TODO:
func NewBolt(dir string) *bolt.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	db, err := bolt.Open(path.Join(dir, "bolt.db"), 0600, nil)
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/etcd-io/bbolt
// - nested buckets
// - transactions
// @TODO:
func NewBBolt(dir string) *bbolt.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	db, err := bbolt.Open(path.Join(dir, "bbolt.db"), 0600, nil)
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/flower-corp/rosedb
// @TODO:
func NewRoseDB(dir string) *rosedb.RoseDB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	opts := rosedb.DefaultOptions(dir)
	opts.LogFileSizeThreshold = 1 << 20

	db, err := rosedb.Open(opts)
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/tidwall/buntdb
// @TODO:
func NewBuntDB(dir string) *buntdb.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	db, err := buntdb.Open(path.Join(dir, "data.db"))
	db.ReadConfig(&buntdb.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/peterbourgon/diskv
// @TODO:
func NewDiskV(dir string) *diskv.Diskv {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	flatTransform := func(s string) []string { return []string{} }

	// Initialize a new diskv store, rooted at "my-data-dir", with a 1MB cache.
	return diskv.New(diskv.Options{
		BasePath:     dir,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
}

// https://github.com/syndtr/goleveldb
// @TODO:
func NewLevelDB(dir string) *leveldb.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/nutsdb/nutsdb
// @TODO:
func NewNutsDB(dir string) *nutsdb.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	opts := nutsdb.DefaultOptions
	opts.Dir = dir

	db, err := nutsdb.Open(opts)
	if err != nil {
		panic(err)
	}

	return db
}

// https://github.com/mattn/go-sqlite3
func NewSQLite(dir string) *sql.DB {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	db, err := sql.Open("sqlite3", path.Join(dir, "data.db"))
	if err != nil {
		panic(err)
	}

	return db
}
