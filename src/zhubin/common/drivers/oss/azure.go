package oss

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/storage"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Azure struct {
	client *storage.Container
}

var ErrMsgAzureObjectNotFound = "storage: service returned error: StatusCode=404, ErrorCode=BlobNotFound, ErrorMessage=The specified blob does not exist."

func NewAzureClient(host, apiVersion, accountName, accountKey string, bucket string) *Azure {
	if !validate(host, apiVersion, accountName, accountKey, bucket) {
		panic("Uninitialized oss driver.")
	}
	client, err := storage.NewClient(accountName, accountKey, host, apiVersion, true)
	if err != nil {
		panic(err)
	}
	blobStorageCLient := client.GetBlobService()
	container := blobStorageCLient.GetContainerReference(bucket)
	exist, err := container.Exists()
	if err != nil {
		panic(err)
	}
	if !exist {
		if err := container.Create(nil); err != nil {
			panic(err)
		}
	}
	return &Azure{client: container}
}

func (azure *Azure) PutObject(tenant string, objKey string, object io.Reader, overwrite bool) error {
	if !overwrite {
		exist, err := azure.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	return azure.client.GetBlobReference(tenant+"/"+objKey).CreateBlockBlobFromReader(object, nil)
}

func (azure *Azure) FPutObject(tenant string, objKey string, filePath string, overwrite bool) error {
	if !overwrite {
		exist, err := azure.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	return azure.client.GetBlobReference(tenant+"/"+objKey).CreateBlockBlobFromReader(file, nil)
}

func (azure *Azure) BPutObject(tenant string, objKey string, objData []byte, overwrite bool) error {
	if !overwrite {
		exist, err := azure.CheckExist(tenant, objKey)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	return azure.client.GetBlobReference(tenant+"/"+objKey).CreateBlockBlobFromReader(bytes.NewReader(objData), nil)
}

func (azure *Azure) GetObject(tenant string, objKey string) (io.ReadCloser, error) {
	obj, err := azure.client.GetBlobReference(tenant + "/" + objKey).Get(nil)
	if err != nil && strings.Contains(err.Error(), ErrMsgAzureObjectNotFound) {
		return nil, ErrNotFound
	}
	return obj, err
}

func (azure *Azure) RemoveObject(tenant string, objKey string) error {
	err := azure.client.GetBlobReference(tenant + "/" + objKey).Delete(nil)
	if err != nil && strings.Contains(err.Error(), ErrMsgAzureObjectNotFound) {
		return nil
	}
	return err
}

func (azure *Azure) ReadObject(tenant string, objKey string) ([]byte, error) {
	obj, err := azure.client.GetBlobReference(tenant + "/" + objKey).Get(nil)
	if err != nil {
		if strings.Contains(err.Error(), ErrMsgAzureObjectNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ioutil.ReadAll(obj)
}

func (azure *Azure) CheckExist(tenant string, objKey string) (bool, error) {
	return azure.client.GetBlobReference(tenant + "/" + objKey).Exists()
}
