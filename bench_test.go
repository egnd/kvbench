package kvbench_test

import (
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
)

var itemsCnt = 10000
var bucketName = []byte("testbuck")

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
