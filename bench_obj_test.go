package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"
)

func Benchmark_ObjBox(b *testing.B) {
	folder := ".data/objbox"

	obx := kvbench.NewObjBox(folder)
	defer obx.Close()

	box := kvbench.BoxForTestObj(obx)

	for i := 0; i < itemsCnt; i++ {
		if _, err := box.Put(kvbench.NewTestObj()); err != nil {
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
			if item.ID, err = box.Put(item); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := box.Get(item.ID); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := box.Query(kvbench.TestObj_.UID.Equals(item.UID, true)).Find(); err != nil {
				b.Error(err)
			}
		}
	})
}
