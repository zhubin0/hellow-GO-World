package oss

import (
	"bytes"
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"io/ioutil"
	"strings"
)

// Compatible with Amazon S3
type Minio struct {
	Bucket   string
	Location string
	client   *minio.Client
}

var ErrMsgMinioObjectNotExists = "The specified key does not exist."

func NewMinioClient(host, port, accessKeyId, SecretAccressKey string, useSSL bool, bucket string, location string) *Minio {
	if !validate(host, port, accessKeyId, SecretAccressKey, bucket) {
		panic("Uninitialized oss driver.")
	}
	endpoint := fmt.Sprintf("%s:%s", host, port)
	client, err := minio.New(endpoint, accessKeyId, SecretAccressKey, useSSL)
	if err != nil {
		panic(err)
	}
	exist, err := client.BucketExists(bucket)
	if err != nil {
		panic(err)
	}
	if !exist {
		if err := client.MakeBucket(bucket, location); err != nil {
			panic(err)
		}
	}
	return &Minio{Bucket: bucket, client: client}
}

func (m *Minio) PutObject(tenant string, objKey string, object io.Reader, overwrite bool) error {
	if !overwrite {
		exist, err := m.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	_, err := m.client.PutObject(m.Bucket, tenant+"/"+objKey, object, "")
	return err
}

func (m *Minio) FPutObject(tenant string, objKey string, filePath string, overwrite bool) error {
	if !overwrite {
		exist, err := m.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	_, err := m.client.FPutObject(m.Bucket, tenant+"/"+objKey, filePath, "")
	return err
}

func (m *Minio) BPutObject(tenant string, objKey string, objData []byte, overwrite bool) error {
	if !overwrite {
		exist, err := m.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	return m.PutObject(tenant, objKey, bytes.NewReader(objData), overwrite)
}

func (m *Minio) GetObject(tenant string, objKey string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(m.Bucket, tenant+"/"+objKey)
	if err != nil && err.Error() == ErrMsgMinioObjectNotExists {
		return nil, ErrNotFound
	}
	return obj, err
}

func (m *Minio) RemoveObject(tenant string, objKey string) error {
	err := m.client.RemoveObject(m.Bucket, tenant+"/"+objKey)
	if err != nil && err.Error() == ErrMsgMinioObjectNotExists {
		return nil
	}
	return err
}

func (m *Minio) ReadObject(tenant string, objKey string) ([]byte, error) {
	obj, err := m.client.GetObject(m.Bucket, tenant+"/"+objKey)
	if err != nil {
		if strings.Contains(err.Error(), ErrMsgMinioObjectNotExists) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ioutil.ReadAll(obj)
}

func (m *Minio) CheckExist(tenant string, objKey string) (bool, error) {
	_, err := m.client.StatObject(m.Bucket, tenant+"/"+objKey)
	if err != nil {
		if err.Error() == ErrMsgMinioObjectNotExists {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
