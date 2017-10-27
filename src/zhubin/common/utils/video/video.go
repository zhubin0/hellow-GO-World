package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"errors"

	"bytes"
	"io"

	"code.google.com/p/log4go"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/afero"
)

var ERR_GPU_SESSION_LIMIT = errors.New("Reach GPU session limit.")

// transcode a video, and return the absolute path of the result video file.
// One should limit the maximum concurrently transcoding sessions to 2, otherwise, error ERR_GPU_SESSION_LIMIT and you should retry later.
func ConvertVideo(inputFile string, workDir string, outputFilePrefix string, useGPU bool, threadPerTask int, videoRateInKb int, videoScale string) (string, error) {
	cmd := onePass(inputFile, workDir, outputFilePrefix, useGPU, threadPerTask, videoRateInKb, videoScale)
	log4go.Debug(cmd)
	outputFilepath := filepath.Join(workDir, outputFilePrefix+".mp4")
	if err := execFfmpegCmd(cmd); err != nil {
		log4go.Error(fmt.Sprintf("Error during video encoding, error: %v", err.Error()))
		// delete the possible failed temp file
		os.Remove(outputFilepath)
		// check the error cause
		if strings.Contains(err.Error(), "OpenEncodeSessionEx failed: out of memory") {
			return "", ERR_GPU_SESSION_LIMIT
		}
		return "", err
	}

	return outputFilepath, nil
}

// take screenshots for a video.
// return the sorted screenshots file names.
func TakeScreenshots(inputFile string, workDir string, outputFilePrefix string, interval int) ([]string, error) {
	cmd := fmt.Sprintf(`ffmpeg -i "%v" -vf fps=1/%v "%v"`, inputFile, interval, fmt.Sprintf("%v_%%d.jpg", getAbsPath(workDir, outputFilePrefix)))
	log4go.Debug(cmd)
	if err := execFfmpegCmd(cmd); err != nil {
		log4go.Info("ffmpeg command: %s", cmd)
		log4go.Error(fmt.Sprintf("Error during screenshot, error: %v", err.Error()))
		return nil, err
	}
	// list the screenshots
	files, err := filterFiles(workDir, outputFilePrefix, []string{".jpg"})
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	return files, nil
}

// create HLS stream from a file.
// NOTE we do not re-encoding the videos and audios, so do make sure the input video is in H.264 format and audio in AAC, MP3, AC-3 or EC-3.
func SegmentHLS(inputFile string, workDir string, outputFilePrefix string, hlsTime int) ([]string, error) {
	if err := execFfmpegCmd(seg_m3u8_cmd(inputFile, workDir, outputFilePrefix, hlsTime)); err != nil {
		log4go.Error(fmt.Sprintf("Error during segment HLS, error: %v", err.Error()))
		return nil, err
	}
	// list the screenshots
	files, err := filterFiles(workDir, outputFilePrefix, []string{".m3u8", ".ts"})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// Video transcoding using one pass.
// e.g. videoScale: "-1:720"
func onePass(inputFile string, workDir string, outputFilePrefix string, useGPU bool, threadPerTask int, videoRateInKb int, videoScale string) string {
	commandName := "ffmpeg"
	buffsize := videoRateInKb * 4
	audioRateInKb := 128
	// use GPU?
	if useGPU {
		return fmt.Sprintf(
			`%v -y -i "%v" -codec:v h264_nvenc -profile:v high -preset slow -b:v %vk -bufsize %vk -vf scale=%v -codec:a aac -b:a %vk -f mp4 "%v"`,
			commandName, inputFile, videoRateInKb, buffsize, videoScale, audioRateInKb, getAbsPath(workDir, fmt.Sprintf("%v.mp4", outputFilePrefix)))
	} else {
		// decide which codecs to use according to configurations
		return fmt.Sprintf(
			`%v -y -i "%v" -codec:v libx264 -profile:v high -preset slow -b:v %vk -bufsize %vk -vf scale=%v -threads %v -codec:a aac -b:a %vk -f mp4 "%v"`,
			commandName, inputFile, videoRateInKb, buffsize, videoScale, threadPerTask, audioRateInKb, getAbsPath(workDir, fmt.Sprintf("%v.mp4", outputFilePrefix)))
	}
}

// seg hls file with original codecs.
func seg_m3u8_cmd(inputFile string, workDir string, outputFilePrefix string, hlsTime int) string {
	commandName := "ffmpeg"
	return fmt.Sprintf(
		`%v -y -i "%v" -codec copy -map 0 -bsf:v h264_mp4toannexb -f segment -segment_list "%v" -segment_format mpegts -segment_time %v -segment_list_type m3u8 "%v"`,
		commandName, inputFile, fmt.Sprintf("%v.m3u8", getAbsPath(workDir, outputFilePrefix)), hlsTime, fmt.Sprintf("%v%%d.ts", getAbsPath(workDir, outputFilePrefix)))
}

// execute a ffmpeg command.
func execFfmpegCmd(fullCommand string) error {
	//we need to split up the command for os.exec
	//parts := strings.Fields(fullCommand)
	parts, err := shellquote.Split(strings.Replace(fullCommand, "\u202A", "", -1)) // fix bug on windows: http://www.fileformat.info/info/unicode/char/202a/index.htm
	if err != nil {
		return err
	}
	head, parts := parts[0], parts[1:]
	cmd := exec.Command(head, parts...)
	cmd.Stdout = os.Stdout
	errbuf := &bytes.Buffer{}
	stderrpipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	go io.Copy(errbuf, stderrpipe)
	if err := cmd.Wait(); err != nil {
		return errors.New(fmt.Sprintf("FFMpeg error: %s, Stderr: %s", err.Error(), string(errbuf.Bytes())))
	}
	return nil
}

// get the absolute file path in the working dir
func getAbsPath(workDir string, filename string) string {
	return filepath.Join(workDir, filename)
}

func filterFiles(root string, prefix string, suffix []string) ([]string, error) {
	type tp struct {
		FileName string
		ModTime  int64
	}
	matched := []tp{}
	var osfs afero.Fs = afero.NewOsFs()
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}
		name := info.Name()
		if strings.HasPrefix(name, prefix) {
			for _, v := range suffix {
				if strings.HasSuffix(name, v) {
					matched = append(matched, tp{FileName: path, ModTime: info.ModTime().UnixNano()})
				}
			}
		}
		return nil
	}
	err := afero.Walk(osfs, root, walkFn)
	if err != nil {
		return nil, err
	}
	// sort
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ModTime < matched[j].ModTime
	})
	ret := []string{}
	for _, v := range matched {
		ret = append(ret, v.FileName)
	}
	return ret, nil
}

// Determines whether we can use GPU or not.
// NOTE if this say false, we won't allow user to enable GPU encoding even if the user force true in configuration file.
func DetectGPU() bool {
	// the command to test if GPU encoding is supported
	// NOTE we do not detect GPU or CUDA themselves, but use the ffmpeg to run an actual null stream. This method is
	//      most reliable.
	// see: https://trac.ffmpeg.org/wiki/Null
	cmd := "ffmpeg -f lavfi -i nullsrc=s=1280x1280:d=1  -map 0:v:0 -c:v h264_nvenc -f null -"
	if err := execFfmpegCmd(cmd); err != nil {
		// any error indicates we can not use GPU or even the FFMPEG itself
		return false
	}
	return true
}
