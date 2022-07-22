package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"

	"github.com/xujiajun/nutsdb"
)

func Benchmark_NutsDB(b *testing.B) {
	folder := ".data/nutsdb"

	db := kvbench.NewNutsDB(folder)
	defer db.Close()

	bucketName := string(bucketName)

	db.Update(func(txn *nutsdb.Tx) error {
		for i := 0; i < itemsCnt; i++ {
			item := kvbench.NewTestObj()
			itemData, _ := json.Marshal(item)
			if err := txn.Put(bucketName, []byte(item.UID), itemData, nutsdb.Persistent); err != nil {
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
			if err := db.Update(func(txn *nutsdb.Tx) error {
				itemData, _ := json.Marshal(item)
				return txn.Put(bucketName, []byte(item.UID), itemData, nutsdb.Persistent)
			}); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			if err := db.View(func(txn *nutsdb.Tx) error {
				entry, err := txn.Get(bucketName, []byte(item.UID))
				if err != nil {
					b.Error(err)
				}

				return json.Unmarshal(entry.Value, &citem)
			}); err != nil {
				b.Error(err)
			}
		}
	})
}
