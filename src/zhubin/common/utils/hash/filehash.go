package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
)

type FileHash struct {
	Sha1   hash.Hash
	Sha256 hash.Hash
}

func NewFileHash() *FileHash {
	return &FileHash{
		Sha1:   sha1.New(),
		Sha256: sha256.New(),
	}
}

//update file hash and return block hash.
func (fh *FileHash) Update(block []byte) (blockHash string, err error) {
	_, err = fh.Sha1.Write(block)
	if err != nil {
		return
	}
	_, err = fh.Sha256.Write(block)
	if err != nil {
		return
	}
	sha1H := sha1.New()
	_, err = sha1H.Write(block)
	if err != nil {
		return
	}
	sha256H := sha256.New()
	_, err = sha256H.Write(block)
	if err != nil {
		return
	}
	blockHash = hex.EncodeToString(sha1H.Sum(nil)) + hex.EncodeToString(sha256H.Sum(nil))
	return
}

//return sha1+sha256 hash
func (fh *FileHash) Sum(b []byte) string {
	return hex.EncodeToString(fh.Sha1.Sum(b)) + hex.EncodeToString(fh.Sha256.Sum(b))
}

//Calculate file hash.
func CalcFileHash(file io.Reader) (string, error) {
	h1 := sha1.New()
	h256 := sha256.New()
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		_, err = h1.Write(buf)
		if err != nil {
			return "", err
		}
		_, err = h256.Write(buf)
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(h1.Sum(nil)) + hex.EncodeToString(h256.Sum(nil)), nil
}

func CalcFileHashAndMime(file io.Reader) (string, string, error) {
	h1 := sha1.New()
	h256 := sha256.New()
	buf := make([]byte, 4096)
	contentType := ""
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", "", err
		}
		if contentType == "" {
			contentType = http.DetectContentType(buf)
		}
		if n == 0 {
			break
		}
		_, err = h1.Write(buf)
		if err != nil {
			return "", "", err
		}
		_, err = h256.Write(buf)
		if err != nil {
			return "", "", err
		}
	}
	return hex.EncodeToString(h1.Sum(nil)) + hex.EncodeToString(h256.Sum(nil)), contentType, nil
}
