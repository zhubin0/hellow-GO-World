package oss

import (
	"datamesh.com/common/drivers/oss/conf"
	"io"
	"strings"
)

var OssClient ObjectStorageDriver

func InitOssDriver(ossConfig *conf.OSS) {
	OssClient = New(ossConfig)
}

// oss storage driver
type ObjectStorageDriver interface {
	// Put an object into oss
	PutObject(tenant string, objKey string, object io.Reader, override bool) error
	// Put an file into oss
	FPutObject(tenant string, objKey string, filePath string, override bool) error
	// Put binary data into oss
	BPutObject(tenant string, objKey string, objData []byte, override bool) error
	// Get an object from oss
	// if object not found, return oss.ErrNotFound as error.
	// NOTE you need to close the stream yourself.
	GetObject(tenant string, objKey string) (io.ReadCloser, error)
	// Read an object from oss to memory.
	// if object not found, return oss.ErrNotFound as error.
	// NOTE you should never use this method when the object is huge.
	ReadObject(tenant string, objKey string) ([]byte, error)
	// Remove an object from oss.
	// If not exists, treat it as deleted.
	RemoveObject(tenant string, objKey string) error
	// Check if the object exists in OSS.
	// return value is valid if error is nil
	CheckExist(tenant string, objKey string) (bool, error)
}

//osType could be:ambry,ali...
func New(ossConfig *conf.OSS) ObjectStorageDriver {
	switch ossConfig.Platform {
	case "Aliyun":
		return NewAliYunClient(
			ossConfig.Host,
			ossConfig.Port,
			ossConfig.AccessKeyID,
			ossConfig.SecretAccessKey,
			ossConfig.Bucket,
		)
	case "Minio", "S3":
		return NewMinioClient(
			ossConfig.Host,
			ossConfig.Port,
			ossConfig.AccessKeyID,
			ossConfig.SecretAccessKey,
			ossConfig.UseSSL,
			ossConfig.Bucket,
			ossConfig.Location,
		)
	case "Azure":
		return NewAzureClient(
			ossConfig.Host,
			ossConfig.ApiVersion,
			ossConfig.AccessKeyID,
			ossConfig.SecretAccessKey,
			ossConfig.Bucket,
		)

	default:
		panic("Unsupported platform: " + ossConfig.Platform)
	}
	return nil
}

func validate(fields ...string) bool {
	for _, v := range fields {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}
