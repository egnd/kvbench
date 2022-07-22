package kvbench_test

import (
	"log"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/egnd/kvbench"
)

func Benchmark_SQLite(b *testing.B) {
	folder := ".data/sqlite"

	db := kvbench.NewSQLite(folder)
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE testtable (
		id INTEGER PRIMARY KEY,
		uid TEXT NOT NULL,
		fname TEXT NOT NULL,
		lname TEXT NOT NULL,
		updated INTEGER NOT NULL
	);`)
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO testtable (uid, fname, lname, updated) values(?,?,?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i := 0; i < itemsCnt; i++ {
		item := kvbench.NewTestObj()
		if _, err := stmt.Exec(item.UID, item.FirstName, item.LastName, item.Updated); err != nil {
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
			if _, err := stmt.Exec(item.UID, item.FirstName, item.LastName, item.Updated); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		var citem kvbench.TestObj
		for i := 0; i < b.N; i++ {
			rows, err := db.Query("SELECT uid, fname, lname, updated FROM testtable WHERE uid=? LIMIT 1", item.UID)
			if err != nil {
				b.Error(err)
			}
			for rows.Next() {
				if err = rows.Scan(&citem.UID, &citem.FirstName, &citem.LastName, &citem.Updated); err != nil {
					b.Error(err)
				}
			}
			rows.Close()
		}
	})
}
