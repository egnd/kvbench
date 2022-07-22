package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"

	"go.etcd.io/bbolt"
)

func Benchmark_BBolt(b *testing.B) {
	folder := ".data/bbolt"

	db := kvbench.NewBBolt(folder)
	defer db.Close()

	db.Update(func(txn *bbolt.Tx) error {
		bucket, _ := txn.CreateBucketIfNotExists(bucketName)
		for i := 0; i < itemsCnt; i++ {
			item := kvbench.NewTestObj()
			itemData, _ := json.Marshal(item)
			if err := bucket.Put([]byte(item.UID), itemData); err != nil {
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
			if err := db.Update(func(txn *bbolt.Tx) error {
				itemData, _ := json.Marshal(item)
				return txn.Bucket(bucketName).Put([]byte(item.UID), itemData)
			}); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			if err := db.View(func(txn *bbolt.Tx) error {
				return json.Unmarshal(txn.Bucket(bucketName).Get([]byte(item.UID)), &citem)
			}); err != nil {
				b.Error(err)
			}
		}
	})
}
