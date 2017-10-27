package mysql

// MySQLConnInfo stores the parameters for MySQL connection.
// Note: currently we do not check the validity of any parameter.
type MySQLConnInfo struct {
	Tag     string `yaml:"tag" bson:"tag"`
	Uname   string `yaml:"uname" bson:"uname"`
	Passwd  string `yaml:"passwd" bson:"passwd"`
	Host    string `yaml:"host" bson:"host"`
	Port    string `yaml:"port" bson:"port"`
	Dbname  string `yaml:"dbname" bson:"dbname"`
	Charset string `yaml:"charset" bson:"charset"`
}
