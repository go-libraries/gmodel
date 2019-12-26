package generator

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestMysqlGenFile(t *testing.T) {

	mysqlHost := "127.0.0.1"
	mysqlPort := "3306"
	mysqlUser := "root"
	mysqlPassword := "sa"
	mysqlDbname := "mgtj"

	dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDbname + "?charset=utf8mb4"

	Mysql := GetMysqlToGo()
	Mysql.Driver.SetDsn(dsn)
	//Mysql.SetStyle("bee")
	Mysql.SetStyle("gorm")
	Mysql.SetModelPath("/Users/limars/Go/src/github.com/go-libraries/gmodel/models")
	Mysql.SetIgnoreTables("cate")
	Mysql.SetPackageName("models")
	Mysql.Run()
}



