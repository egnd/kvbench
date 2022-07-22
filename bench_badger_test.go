package kvbench_test

import (
	"log"
	"testing"

	"github.com/egnd/kvbench"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/dustin/go-humanize"
)

func Benchmark_Badger(b *testing.B) {
	folder := ".data/badger"

	db := kvbench.NewBadgerDB(folder)
	defer db.Close()

	db.Update(func(txn *badger.Txn) error {
		for i := 0; i < itemsCnt; i++ {
			item := kvbench.NewTestObj()
			itemData, _ := json.Marshal(item)
			if err := txn.Set([]byte(item.UID), itemData); err != nil {
				panic(err)
			}
		}
		return nil
	})

	size, err := DirSize(folder)
	if err != nil {
		panic(err)
	}

	b.Cleanup(func() { log.Println("-------- size:", itemsCnt, humanize.Bytes(uint64(size))) })

	item := kvbench.NewTestObj()

	b.Run("put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := db.Update(func(txn *badger.Txn) error {
				itemData, _ := json.Marshal(item)
				return txn.Set([]byte(item.UID), itemData)
			}); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			if err := db.View(func(txn *badger.Txn) error {
				res, err := txn.Get([]byte(item.UID))
				if err != nil {
					return nil
				}

				return res.Value(func(val []byte) error {
					return json.Unmarshal(val, &citem)
				})
			}); err != nil {
				b.Error(err)
			}
		}
	})
}
