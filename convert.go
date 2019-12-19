/*
 * Copyright (c) 2019 Mars Lee. All rights reserved.
 */

package model

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type SqlDriver interface {
	SetDsn(dsn string, options ...interface{})
	Connect() error
	ReadTablesColumns(table string) []Column
	GetTables() []string
}

type Convert struct {
	ModelPath   string	// save path
	DriverType  string	// driver name like mysql postgre_sql sql_server ......
	TagKey		string  // tab key save like gorm orm ......
	PackageName string  // go package name

	TablePrefix  map[string]string    //if table exists prefix
	TableColumn  map[string][]Column  //key is table , value is Column list
	IgnoreTables []string   // ignore tables
	Tables       []string   // all tables

	Driver       SqlDriver  // impl SqlDriver instance
}

//get real gen tables as []string
func (convert *Convert) getGenTables() []string {
	tables := make([]string, 0)
	convert.Tables = convert.Driver.GetTables()
	for _, table := range convert.Tables {
		isIgnore := false
		for _, ignore := range convert.IgnoreTables {
			if table == ignore {
				isIgnore = true
				break
			}
		}

		if !isIgnore {
			tables = append(tables, table)
		}
	}

	return tables
}

//set table prefix
//if exists
//replace prefix to empty string
func (convert *Convert) SetTablePrefix(table, prefix string) {
	convert.TablePrefix[table] = prefix
}

// set model save path
func (convert *Convert) SetModelPath(path string) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("path not exists with error：%v", err))
		}
		log.Println(fmt.Sprintf("path error：%v", err))
	}

	convert.ModelPath = path
}

// set model save path
func (convert *Convert) SetIgnoreTables(table ...string) {
	convert.IgnoreTables = append(convert.IgnoreTables, table...)
}

// set model save path
func (convert *Convert) SetPackageName(name string) {
	convert.PackageName = name
}

//run
//1. connect
//2. getTable
//3. getColumns
//4. build
//5. write file
func (convert *Convert) Run() {

	err := convert.Driver.Connect()
	if err != nil {
		panic(err)
	}

	for _, tableRealName := range convert.getGenTables() {
		prefix, ok := convert.TablePrefix[tableRealName]
		if ok {
			tableRealName = tableRealName[len(prefix):]
		}
		tableName := tableRealName

		switch len(tableName) {
		case 0:
			continue
		case 1:
			tableName = strings.ToUpper(tableName[0:1])
		default:
			tableName = strings.ToUpper(tableName[0:1]) + tableName[1:]
		}

		columns := convert.Driver.ReadTablesColumns(tableRealName)
		content := convert.build(tableName, tableRealName, prefix, columns)
		convert.writeModel(tableRealName, content) //写文件
	}
}

//build content with table info
func (convert *Convert) build(tableName, tableRealName, prefix string, columns []Column) (content string) {
	depth := 1
	content += "package " + convert.PackageName + "\n\n" //写包名
	content += "type " + tableName + " struct {\n"

	for _, v := range columns {
		var comment string
		if v.ColumnComment != "" {
			comment = fmt.Sprintf(" // %s", v.ColumnComment)
		}
		content += fmt.Sprintf("%s%s %s %s%s\n",
			Tab(depth), v.GetGoColumn(prefix, true), v.GetGoType(), v.GetTag("orm"), comment)
	}
	content += Tab(depth-1) + "}\n\n"

	content += fmt.Sprintf("func (%s *%s) %s() string {\n",
		LcFirst(tableName), tableName, "GetTableName")
	content += fmt.Sprintf("%sreturn \"%s\"\n",
		Tab(depth), tableRealName)
	content += "}\n\n"
	return content
}

//write file
func (convert *Convert) writeModel(name, content string) {
	filePath := fmt.Sprintf("%s/%s.go", convert.ModelPath, name)
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Can not write file" + filePath)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.WriteString(content)
	if err != nil {
		log.Println("Can not write file" + filePath)
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	_ = cmd.Run()
}
