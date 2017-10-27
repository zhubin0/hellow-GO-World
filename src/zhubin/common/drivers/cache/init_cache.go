package cache

var L2_CACHE_CLIENT L2Cache

func InitL2CacheClient(cacheConfig *Cache) {
	if cacheConfig.IsCluster {
		L2_CACHE_CLIENT = NewRedisCluster(cacheConfig.Address, cacheConfig.Password)
		return
	}
	L2_CACHE_CLIENT = NewRedis(cacheConfig.Address, cacheConfig.Password, cacheConfig.DbNum)
}