package hash

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestHash(t *testing.T) {
	fileHash := NewFileHash()

	file, err := os.Open("/Users/smile/Downloads/aaa.jpg")
	if err != nil {
		panic(err)
	}
	i := 0
	buf := make([]byte, 1024*4)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		//TODO to support upload empty file.
		if 0 == n {
			break
		}
		i++
		fmt.Println(i)
		blockHash, err := fileHash.Update(buf[:n])
		if err != nil {
			panic(err)
		}
		t.Log(blockHash)
	}

	fhash := fileHash.Sum(nil)
	t.Log(fhash)
}
