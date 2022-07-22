package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"
)

func Benchmark_LevelDB(b *testing.B) {
	folder := ".data/leveldb"

	db := kvbench.NewLevelDB(folder)
	defer db.Close()

	for i := 0; i < itemsCnt; i++ {
		item := kvbench.NewTestObj()
		itemData, _ := json.Marshal(item)
		if err := db.Put([]byte(item.UID), itemData, nil); err != nil {
			panic(err)
		}
	}

	size, err := DirSize(folder)
	if err != nil {
		panic(err)
	}

	b.Cleanup(func() { log.Println("-------- size:", itemsCnt, humanize.Bytes(uint64(size))) })

	item := kvbench.NewTestObj()

	b.Run("put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			itemData, _ := json.Marshal(item)
			if err := db.Put([]byte(item.UID), itemData, nil); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			data, err := db.Get([]byte(item.UID), nil)
			if err != nil {
				b.Error(err)
			}

			if err := json.Unmarshal(data, &citem); err != nil {
				b.Error(err)
			}
		}
	})
}
