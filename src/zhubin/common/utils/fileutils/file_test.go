package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {

	file, err := os.Open("/Users/smile/Downloads/goroutine.gz")
	if err != nil {
		panic(err)
	}

	fmt.Print("file name is: ")
	fmt.Println(file.Name())

	fmt.Println(filepath.Base(file.Name()))
}

func TestGetParentDirectory(t *testing.T) {
	path := "/Users/nick/Desktop/soft/mygo2/src/datamesh.com/holo-server/portal/packages/SpectatorView_1.0.40.0_x86__pzq3xp76mxafg/config.ini"
	parentDir := GetParentDirectory(path)
	defer os.RemoveAll(parentDir)
}
