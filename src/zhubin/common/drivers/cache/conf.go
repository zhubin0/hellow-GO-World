package cache

type Cache struct {
	Address   string `yaml:"address"` // host:port
	Password  string `yaml:"password"`
	DbNum     int    `yaml:"db_num"` // not applicable on cluster
	IsCluster bool   `yaml:"is_cluster"`
}