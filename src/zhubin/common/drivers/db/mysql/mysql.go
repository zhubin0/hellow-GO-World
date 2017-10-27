package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// defines mysql connection parameters and the client. It implements repository
// interface so it can be pooling.
type MySQLDB struct {
	para *MySQLConnInfo
	Db   gorp.DbMap
}

func NewMySQLDb(conf *MySQLConnInfo) (*MySQLDB, error) {
	logrus.Info(fmt.Sprintf("Initializing: connecting to mysql instance '%s', host: %s, port: %s", conf.Dbname, conf.Host, conf.Port))
	// init single mysql instance
	m := &MySQLDB{para: conf}
	if err := m.init(); err != nil {
		return nil, err
	}
	return m, nil
}

// ping the underlying database.
func (m *MySQLDB) Ping() error {
	return m.Db.Db.Ping()
}

// generate the connection URI.
func (m *MySQLDB) uri() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", m.para.Uname, m.para.Passwd, m.para.Host, m.para.Port, m.para.Dbname)
}

// initialize mysql instance against given MySQLConnInfo.
func (m *MySQLDB) init() error {
	// connect to db using standard Go database/sql API
	db, err := sql.Open("mysql", m.uri())
	if err != nil {
		return fmt.Errorf("failed to connect MySQL database %s, message: %v", m.para.Host, err)
	}
	// construct a gorp DbMap
	m.Db = gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", m.para.Charset}}
	m.Db.TraceOn("[gorp]", logrus.StandardLogger())

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	//dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	//err = dbmap.CreateTablesIfNotExists()
	//checkErr(err, "Create tables failed")

	return nil
}
