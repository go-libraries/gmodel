package model

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var TypeMappingMysqlToGo = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "string", // time.Time
	"datetime":           "string", // time.Time
	"timestamp":          "string", // time.Time
	"time":               "string", // time.Time
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}
var tableToGo *Convert
var syncMysql sync.Once

type MysqlToGo struct {
	Dsn string
	db  *sql.DB
}

func GetMysqlToGo() *Convert {

	syncMysql.Do(func() {
		tableToGo = &Convert{
			ModelPath:    "",
			TablePrefix:  make(map[string]string),
			TableColumn:  make(map[string][]Column),
			IgnoreTables: make([]string, 0),
			Tables:       make([]string, 0),
			DriverType:   "mysql",
			Driver: &MysqlToGo{
				Dsn: "",
				db:  nil,
			},
		}
	})
	return tableToGo
}

//connection to mysql
func (mtg *MysqlToGo) SetDsn(dsn string, options ...interface{}) {
	mtg.Dsn = dsn
}

//connection to mysql
func (mtg *MysqlToGo) Connect() error {
	db, err := sql.Open("mysql", mtg.Dsn)
	if err != nil {
		return err
	}

	mtg.db = db
	return nil
}

// tables
func (mtg *MysqlToGo) GetTables() (tables []string) {
	rows, err := mtg.db.Query("show tables;")
	if err != nil {
		return tables
	}

	if rows == nil {
		return tables
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var f string
		err := rows.Scan(&f)
		if err != nil {
			panic(err)
		}
		tables = append(tables, f)
	}

	return tables
}

//read struct from db
func (mtg *MysqlToGo) ReadTablesColumns(table string) []Column {
	columns := make([]Column, 0)
	rows, err := mtg.db.Query(fmt.Sprintf(`SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE()  AND TABLE_NAME = '%s'`, table))

	if err != nil {
		log.Printf("table rows is nil with table:%s error: %v \n", table, err)
		return columns
	}

	if rows == nil {
		log.Printf("rows is nil with table:%s \n", table)
		return columns
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {

		col := Column{}
		err = rows.Scan(&col.ColumnName, &col.Type, &col.Nullable, &col.TableName, &col.ColumnComment)
		col.Tag = col.ColumnName

		if err != nil {
			log.Println(err.Error())
			continue
		}

		columns = append(columns, col)
	}
	return columns
}
