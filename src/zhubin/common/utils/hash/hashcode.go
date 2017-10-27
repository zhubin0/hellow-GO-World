package hash

import (
	"fmt"
	"hash/fnv"
)

func FNV1aHashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

const HLS_SUB_DIR_COUNT = 1000

func CalHLSSubDir(videoPrefix string) string {
	return fmt.Sprintf("s%d", HLS_SUB_DIR_COUNT+(FNV1aHashCode(videoPrefix)%HLS_SUB_DIR_COUNT))
}

func ListHLSSubDirs() (list []string) {
	for i := HLS_SUB_DIR_COUNT; i < 2*HLS_SUB_DIR_COUNT-1; i++ {
		list = append(list, fmt.Sprintf("s%d", i))
	}
	return
}
