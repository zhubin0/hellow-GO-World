package oss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"strings"
)

type AliYun struct {
	client *oss.Bucket
}

var ErrMsgAliyunObjectNotFound = "oss: service returned error: StatusCode=404, ErrorCode=NoSuchKey, ErrorMessage=The specified key does not exist."

func NewAliYunClient(host, port, accessKeyId, SecretAccressKey, bucket string) *AliYun {
	if !validate(host, port, accessKeyId, SecretAccressKey, bucket) {
		panic("Uninitialized oss driver.")
	}

	endpoint := fmt.Sprintf("%s:%s", host, port)

	client, err := oss.New(endpoint, accessKeyId, SecretAccressKey)
	if err != nil {
		panic(err)
	}
	exist, err := client.IsBucketExist(bucket)
	if err != nil {
		panic(err)
	}
	if !exist {
		if err := client.CreateBucket(bucket); err != nil {
			panic(err)
		}
	}

	bkt, err := client.Bucket(bucket)
	if err != nil {
		panic(err)
	}
	return &AliYun{client: bkt}
}

func (ali *AliYun) PutObject(tenant string, objKey string, object io.Reader, overwrite bool) error {
	if !overwrite {
		exist, err := ali.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	return ali.client.PutObject(tenant+"/"+objKey, object)
}

func (ali *AliYun) FPutObject(tenant string, objKey string, filePath string, overwrite bool) error {
	if !overwrite {
		exist, err := ali.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	return ali.client.PutObjectFromFile(tenant+"/"+objKey, filePath)
}

func (ali *AliYun) BPutObject(tenant string, objKey string, objData []byte, overwrite bool) error {
	if !overwrite {
		exist, err := ali.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	return ali.client.PutObject(tenant+"/"+objKey, bytes.NewReader(objData))
}

func (ali *AliYun) GetObject(tenant string, objKey string) (io.ReadCloser, error) {
	obj, err := ali.client.GetObject(tenant + "/" + objKey)
	if err != nil && strings.Contains(err.Error(), ErrMsgAliyunObjectNotFound) {
		return nil, ErrNotFound
	}
	return obj, err
}

func (ali *AliYun) RemoveObject(tenant string, objKey string) error {
	err := ali.client.DeleteObject(tenant + "/" + objKey)
	if err != nil && strings.Contains(err.Error(), ErrMsgAliyunObjectNotFound) {
		return nil
	}
	return err
}

func (ali *AliYun) ReadObject(tenant string, objKey string) ([]byte, error) {
	obj, err := ali.client.GetObject(tenant + "/" + objKey)
	if err != nil {
		if strings.Contains(err.Error(), ErrMsgAliyunObjectNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ioutil.ReadAll(obj)
}

func (ali *AliYun) CheckExist(tenant string, objKey string) (bool, error) {
	return ali.client.IsObjectExist(tenant + "/" + objKey)
}
