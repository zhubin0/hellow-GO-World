package cache

import (
	"datamesh.com/common/utils"
	"datamesh.com/common/utils/randgen"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	L2_CACHE_CLIENT := NewRedis("Redis", "127.0.0.1:6379", 0)
	defer L2_CACHE_CLIENT.Close()

	key := randgen.GenUniqueString(8)
	data := randgen.GenRandNumString(64)

	err := L2_CACHE_CLIENT.Save(key, data, 1)
	assert.Nil(t, err)
	ret, err := L2_CACHE_CLIENT.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, data, ret)

	// testing delete
	exist, err := L2_CACHE_CLIENT.Exists(key)
	assert.Nil(t, err)
	assert.True(t, exist)
	err = L2_CACHE_CLIENT.Delete(key)
	assert.Nil(t, err)
	exist, err = L2_CACHE_CLIENT.Exists(key)
	assert.Nil(t, err)
	assert.False(t, exist)

	// testing expire
	err = L2_CACHE_CLIENT.Save(key, data, 1)
	assert.Nil(t, err)
	exist, err = L2_CACHE_CLIENT.Exists(key)
	assert.Nil(t, err)
	assert.True(t, exist)
	time.Sleep(time.Second * 2)
	exist, err = L2_CACHE_CLIENT.Exists(key)
	assert.Nil(t, err)
	assert.False(t, exist)

	// testing incr
	err = L2_CACHE_CLIENT.Save(key, "1", 1)
	assert.Nil(t, err)
	_, err = L2_CACHE_CLIENT.Incr(key)
	assert.Nil(t, err)
	ret, err = L2_CACHE_CLIENT.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, "2", ret)
}

func TestRedisCluster(t *testing.T) {
	cluster := NewRedisCluster("100.64.10.1:6379", "")

	var round int64 = 100000
	var i int64
	prefix := "test_"
	t1 := time.Now().UnixNano()
	for i = 0; i < round; i++ {
		if err := cluster.Save(prefix+randgen.GenRandString(32), randgen.GenUniqueString(32), 1000); err != nil {
			panic(err)
		}
	}
	t2 := time.Now().UnixNano()
	fmt.Printf("Redis Cluster Performance: write QPS %d\n", 1e9*round/(t2-t1))
}