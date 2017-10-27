package oss

import (
	"datamesh.com/MeshExpert/config"
	"datamesh.com/common/drivers/oss/conf"
	"datamesh.com/common/utils/randgen"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	bucket = "testbucket"
	tenant = "19032010"
)

func initConfig(platform, host, port, accessKeyID, secretAccressKey, bucket string, useSSL bool) {
	config.CONFIGS = &config.MeshExpertConfigs{
		OSS: &conf.OSS{
			Platform:        platform,
			Host:            host,
			Port:            port,
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccressKey,
			Bucket:          bucket,
			UseSSL:          useSSL,
			ApiVersion:      "2016-05-31",
			Location:        "minio",
		},
	}
}

func TestMinio_FPutObject(t *testing.T) {
	initConfig("Minio", "192.168.1.52", "9000", "KVQVTQKMM34JATV3G327", "8Vh7fS2eNc0pVqYCzv2IqWf3kP+arLfDw3vwAcrz", bucket, false)
	InitOssDriver(config.CONFIGS.OSS)

	objKey := "aaaaaaa"
	reader := strings.NewReader("aaaaaaaa")
	err := OssClient.PutObject("ttt", objKey, reader, true)
	assert.Nil(t, err)
}

func TestMinio_CheckExist(t *testing.T) {
	initConfig("Minio", "192.168.1.52", "9000", "KVQVTQKMM34JATV3G327", "8Vh7fS2eNc0pVqYCzv2IqWf3kP+arLfDw3vwAcrz", bucket, false)
	InitOssDriver(config.CONFIGS.OSS)

	objKey := "aaaaasaa"
	rd, err := OssClient.CheckExist("ttt", objKey)

	fmt.Println(rd, err)

}

func TestMinio_GetObject(t *testing.T) {
	initConfig("Minio", "192.168.1.52", "9000", "KVQVTQKMM34JATV3G327", "8Vh7fS2eNc0pVqYCzv2IqWf3kP+arLfDw3vwAcrz", bucket, false)
	InitOssDriver(config.CONFIGS.OSS)

	objKey := "aaaaaaa"
	rd, err := OssClient.GetObject("ttt", objKey)

	data, err := ioutil.ReadAll(rd)

	fmt.Println(data)
	fmt.Println(err)
}

func TestMinio(t *testing.T) {
	bucket = "withcheck"
	initConfig("Minio", "192.168.1.52", "9000", "KVQVTQKMM34JATV3G327", "8Vh7fS2eNc0pVqYCzv2IqWf3kP+arLfDw3vwAcrz", bucket, false)
	InitOssDriver(config.CONFIGS.OSS)

	objKey := randgen.GenUniqueString(8)
	reader := strings.NewReader(objKey)
	err := OssClient.PutObject(tenant, objKey, reader, true)
	assert.Nil(t, err)

	ret, err := OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	// overwrite
	strData1 := randgen.GenUniqueString(8)
	reader1 := strings.NewReader(strData1)
	err = OssClient.PutObject(tenant, objKey, reader1, true)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, strData1, string(data))

	// do not overwrite
	strData2 := randgen.GenUniqueString(8)
	reader2 := strings.NewReader(strData2)
	err = OssClient.PutObject(tenant, objKey, reader2, false)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, strData1, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)

	err = OssClient.BPutObject(tenant, objKey, []byte(objKey), true)
	assert.Nil(t, err)
	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)

}

func TestS3(t *testing.T) {
	bucket = "datamesh-test"
	initConfig("S3", "s3-ap-northeast-1.amazonaws.com", "443", "AKIAJ7DT56NHUVKO6OPQ", "WsaeGn7WMU/2XmeCIJievoIzfhffRbEqcK6n9xuN", bucket, true)
	InitOssDriver(config.CONFIGS.OSS)

	objKey := randgen.GenUniqueString(8)
	reader := strings.NewReader(objKey)
	err := OssClient.PutObject(tenant, objKey, reader, true)
	panic(err)
	assert.Nil(t, err)

	ret, err := OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)

	err = OssClient.BPutObject(tenant, objKey, []byte(objKey), true)
	panic(err)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)
}

func TestAliYun_GetObject(t *testing.T) {
	bucket = "me-holocloud-image"
	initConfig("Aliyun", "oss-cn-beijing.aliyuncs.com", "80", "LTAIjrFhcnc3Rg34", "0Bn2h6DIDMpN3LTfbDxlLYNU6GPKYD", bucket, false)
	InitOssDriver(config.CONFIGS.OSS)

	obj, err := OssClient.GetObject("ttt", "ttt")
	fmt.Println(err)
	fmt.Println(obj)
}

func TestAliYun(t *testing.T) {

	bucket = "me-holocloud-image"
	initConfig("Aliyun", "oss-cn-beijing.aliyuncs.com", "80", "LTAIjrFhcnc3Rg34", "0Bn2h6DIDMpN3LTfbDxlLYNU6GPKYD", bucket, false)
	InitOssDriver(config.CONFIGS.OSS)

	objKey := randgen.GenUniqueString(8)
	reader := strings.NewReader(objKey)
	err := OssClient.PutObject(tenant, objKey, reader, true)
	assert.Nil(t, err)

	ret, err := OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	// overwrite
	strData1 := randgen.GenUniqueString(8)
	reader1 := strings.NewReader(strData1)
	err = OssClient.PutObject(tenant, objKey, reader1, true)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, strData1, string(data))

	// do not overwrite
	strData2 := randgen.GenUniqueString(8)
	reader2 := strings.NewReader(strData2)
	err = OssClient.PutObject(tenant, objKey, reader2, false)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, strData1, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)

	err = OssClient.BPutObject(tenant, objKey, []byte(objKey), true)
	assert.Nil(t, err)
	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)

}

func TestAzure_GetObject(t *testing.T) {
	bucket = "datameshtest"
	initConfig("Azure", "core.chinacloudapi.cn", "443", "datameshtest", "8yA6/aI21+fR2K0aktFCcyCbvhzLFCKoOfOQB9rqwJgyloeY31CFWP0zDA42yh5SRctiudjGI06yn6trNKTyNg==", bucket, true)
	InitOssDriver(config.CONFIGS.OSS)

	obj, err := OssClient.GetObject("ttt", "ttt")
	fmt.Println(obj)
	fmt.Println(err)
}

func TestAzure(t *testing.T) {
	bucket = "datameshsocial"
	initConfig("Azure", "core.chinacloudapi.cn", "443", "datameshsocial", "iuJSMsxbeoppvVdig2HXcuBf0AY4DkDaAH0wulOedWxXNCAbmTMhWy5jYrhoEip3Xs7s8oMc1LB03KYISzrCbA==", bucket, true)
	InitOssDriver(config.CONFIGS.OSS)
	objKey := randgen.GenUniqueString(8)
	reader := strings.NewReader(objKey)
	err := OssClient.PutObject(tenant, objKey, reader, true)
	assert.Nil(t, err)

	ret, err := OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	// overwrite
	strData1 := randgen.GenUniqueString(8)
	reader1 := strings.NewReader(strData1)
	err = OssClient.PutObject(tenant, objKey, reader1, true)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, strData1, string(data))

	// do not overwrite
	strData2 := randgen.GenUniqueString(8)
	reader2 := strings.NewReader(strData2)
	err = OssClient.PutObject(tenant, objKey, reader2, false)
	assert.Nil(t, err)

	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, strData1, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)

	err = OssClient.BPutObject(tenant, objKey, []byte(objKey), true)
	assert.Nil(t, err)
	ret, err = OssClient.GetObject(tenant, objKey)
	assert.Nil(t, err)
	data, err = ioutil.ReadAll(ret)
	assert.Nil(t, err)
	assert.Equal(t, objKey, string(data))

	err = OssClient.RemoveObject(tenant, objKey)
	assert.Nil(t, err)
}
