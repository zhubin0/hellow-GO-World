package video

import (
	"testing"

	"datamesh.com/common/utils"
	"github.com/stretchr/testify/assert"
)

func TestDetectGPU(t *testing.T) {
	ok := DetectGPU()
	assert.True(t, ok) // change according to machines, and should run on multiple machines (with or without GPU)
}

// test spaced folders
func TestTakeScreenshots(t *testing.T) {
	videoTmpFile := "‪C:\\Users\\DataMesh\\Desktop\\test folder\\test1.mp4"
	workDir := "E:\\Video\\WorkDir\\test dir"
	_, err := TakeScreenshots(videoTmpFile, workDir, utils.GenRandString(8), 10)
	assert.Nil(t, err)
}
