package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"

	"github.com/tidwall/buntdb"
)

func Benchmark_BuntDB(b *testing.B) {
	folder := ".data/buntdb"

	db := kvbench.NewBuntDB(folder)
	defer db.Close()

	db.Update(func(txn *buntdb.Tx) error {
		for i := 0; i < itemsCnt; i++ {
			item := kvbench.NewTestObj()
			itemData, _ := json.Marshal(item)
			if _, _, err := txn.Set(item.UID, string(itemData), nil); err != nil {
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
			if err := db.Update(func(txn *buntdb.Tx) error {
				itemData, _ := json.Marshal(item)
				if _, _, err := txn.Set(item.UID, string(itemData), nil); err != nil {
					panic(err)
				}
				return err
			}); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			if err := db.View(func(txn *buntdb.Tx) error {
				data, err := txn.Get(item.UID)
				if err != nil {
					return err
				}

				return json.Unmarshal([]byte(data), &citem)
			}); err != nil {
				b.Error(err)
			}
		}
	})
}
