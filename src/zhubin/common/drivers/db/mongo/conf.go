package mongo

type Mongo struct {
	Address  string `yaml:"address"` // host:port or host:port list
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"` /// set by program, not by user
}
