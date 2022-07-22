package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"
)

func Benchmark_DiskV(b *testing.B) {
	folder := ".data/diskv"

	db := kvbench.NewDiskV(folder)

	for i := 0; i < itemsCnt; i++ {
		item := kvbench.NewTestObj()
		itemData, _ := json.Marshal(item)
		if err := db.Write(item.UID, itemData); err != nil {
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
			if err := db.Write(item.UID, itemData); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			data, err := db.Read(item.UID)
			if err != nil {
				b.Error(err)
			}
			if err := json.Unmarshal(data, &citem); err != nil {
				b.Error(err)
			}
		}
	})
}
