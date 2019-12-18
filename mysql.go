package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
)

var tableToGo *MysqlToGo
var syncMysql sync.Once
type MysqlToGo struct {
	*Convert
	Dsn string
	db *sql.DB
}

func GetMysqlToGo() *MysqlToGo  {
	syncMysql.Do(func() {
		tableToGo = &MysqlToGo{
			Dsn:          "",
			db:           nil,
			Convert: &Convert{
				ModelPath:    "",
				TablePrefix:  make(map[string]string),
				TableColumn:  make(map[string][]column),
				IgnoreTables: make([]string, 0),
				Tables:       make([]string, 0),
			},
		}
	})
	return tableToGo
}

//connection to mysql
func (mtg *MysqlToGo) SetDsn(dsn string)  {
	db, err := sql.Open("mysql", dsn)
	if err!= nil{
		panic(err)
	}

	mtg.Dsn = dsn
	mtg.db = db
}

//set table prefix
//if exists
//replace prefix to empty string
func (mtg *MysqlToGo) SetTablePrefix(table,prefix string)  {
	mtg.TablePrefix[table] = prefix
}

// set model save path
func (mtg *MysqlToGo) SetModelPath(path string)  {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			log.Panicf("path not exists with error：%v \n", err)
		}
		log.Printf("path error：%v \n", err)
	}

	mtg.ModelPath = path
}

// set model save path
func (mtg *MysqlToGo) SetIgnoreTables(table... string)  {
	mtg.IgnoreTables = append(mtg.IgnoreTables, table...)
}

func (mtg *MysqlToGo) GetTables() []string  {
	rows, err := mtg.db.Query("show tables;")
	if err != nil {
		panic(err)
	}

	if rows == nil {
		panic("rows is nil")
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
		mtg.Tables = append(mtg.Tables, f)
	}

	return mtg.Tables
}

//read struct from db
func (mtg *MysqlToGo) ReadTablesColumns() {
	for _,table := range mtg.Tables {
		isIgnore := false
		for _,ignore := range mtg.IgnoreTables {
			if table == ignore {
				isIgnore = true
				break
			}
		}

		if !isIgnore {
			mtg.readTablesColumns(table)
		}
	}

}

//read struct from db
func (mtg *MysqlToGo) readTablesColumns(table string) {
	result,err := mtg.db.Query(fmt.Sprintf(`SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE()  AND TABLE_NAME = '%s'`, table))

	if err != nil {
		log.Printf("table result is nil with table:%s error: %v \n",table,err)
		return
	}

	if result == nil {
		log.Printf("result is nil with table:%s \n",table)
		return
	}


	for result.Next() {

		col := column{}
		err = result.Scan(&col.ColumnName, &col.Type, &col.Nullable, &col.TableName, &col.ColumnComment)
		col.Tag = col.ColumnName

		if err != nil {
			log.Println(err.Error())
			continue
		}

		mtg.TableColumn[table] = append(mtg.TableColumn[table], col)
	}
}





