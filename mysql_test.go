package model

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestMysqlGenFile(t *testing.T) {

	mysqlHost := "127.0.0.1"
	mysqlPort := "3306"
	mysqlUser := "root"
	mysqlPassword := "sa"
	mysqlDbname := "blog"

	dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDbname + "?charset=utf8mb4"

	Mysql := GetMysqlToGo()
	Mysql.SetDsn(dsn)
	Mysql.SetModelPath("/tmp")
	Mysql.GetTables()
	Mysql.ReadTablesColumns()
	Mysql.SetPackageName("models")
	Mysql.Run()
}

