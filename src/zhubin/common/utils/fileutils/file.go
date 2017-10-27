package fileutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"archive/zip"
	"compress/gzip"
	"errors"
	"path/filepath"
	"strings"

	"datamesh.com/common/utils/hash"
	"github.com/phayes/permbits"
	"github.com/spf13/afero"
)

type Size interface {
	Size() int64
}
type Stat interface {
	Stat() (os.FileInfo, error)
}

// get the size of a file.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return 0, errors.New("not a file")
	}
	return info.Size(), nil
}

// Get file size.
func FileSize(reader io.Reader) (int64, error) {
	size, ok := reader.(Size)
	if !ok {
		stat, sok := reader.(Stat)
		if !sok {
			return 0, fmt.Errorf("Can not assert reader(type %T) to Size or Stat.", reader)
		}
		fileInfo, err := stat.Stat()
		if err != nil {
			return 0, err
		}
		return fileInfo.Size(), nil
	}
	return size.Size(), nil
}

func IsFileOrDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func MakePrivateDir(path string) error {
	return os.MkdirAll(path, 0700)
}

// determine the type of the path: true for directory, error for not exists
// or other errors
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// does the path represent a regular file
func IsRegularFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), nil
}

// does the path readable and writable
func ReadableWritable(path string) (bool, error) {
	perms, err := permbits.Stat(path)
	if err != nil {
		return false, err
	}
	return perms.UserRead() && perms.UserWrite(), nil
}

//write data to a file.
func WriteFile(filePath string, reader io.Reader) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 1024*1024*4)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		_, err = file.Write(buf[:n])
		if err != nil {
			return nil, err
		}
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return os.Open(filePath)
}

//create a file.
func CreateFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

//read file.
func ReadFile(filePath string) ([]byte, error) {
	exist, err := IsFileOrDirExists(filePath)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("File %s dose not exists.", filePath)
	}
	regular, err := IsRegularFile(filePath)
	if err != nil {
		return nil, err
	}
	if !regular {
		return nil, fmt.Errorf("File %s is not a regular file", filePath)
	}
	return readFile(filePath)
}

// open file for read.
func OpenReadFile(filePath string) (*os.File, error) {
	exist, err := IsFileOrDirExists(filePath)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("File %s dose not exists.", filePath)
	}
	regular, err := IsRegularFile(filePath)
	if err != nil {
		return nil, err
	}
	if !regular {
		return nil, fmt.Errorf("File %s is not a regular file", filePath)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

//compress file use gzip and store it to target directory.
//return filesize,commpressed file size and error.
func GZIP(source, target, filename string) (int64, int64, error) {
	reader, err := os.Open(source)
	if err != nil {
		return 0, 0, err
	}
	fileInfo, err := reader.Stat()
	if err != nil {
		return 0, 0, err
	}
	fileSize := fileInfo.Size()
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return 0, 0, err
	}
	defer writer.Close()
	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()
	n, err := io.Copy(archiver, reader)
	return fileSize, n, err
}

// compress file arrays
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func GetParentDirectory(path string) string {
	runes := []rune(path)
	l := 0 + strings.LastIndex(path, "/")
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[0:l])
}

// make all the paths and their parent paths if not exists
func MakeDirsIfNotExist(paths []string) error {
	for _, p := range paths {
		exist, err := afero.Exists(afero.OsFs{}, p)
		if err != nil {
			return err
		}
		if !exist {
			err := os.MkdirAll(p, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MakeHlsSublocations(rootPath string) error {
	subs := hash.ListHLSSubDirs()
	paths := []string{}
	for _, s := range subs {
		paths = append(paths, filepath.Join(rootPath, s))
	}
	return MakeDirsIfNotExist(paths)
}
