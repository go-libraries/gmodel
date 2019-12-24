/*
 * Copyright (c) 2019 Mars Lee. All rights reserved.
 */

package generator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type SqlDriver interface {
	SetDsn(dsn string, options ...interface{})
	Connect() error
	ReadTablesColumns(table string) []Column
	GetTables() []string
	GetDriverType() string
}

type Convert struct {
	ModelPath   string // save path
	Style       string // tab key save like gorm ,orm ,bee orm......
	PackageName string // go package name

	TablePrefix  map[string]string   //if table exists prefix
	TableColumn  map[string][]Column //key is table , value is Column list
	IgnoreTables []string            // ignore tables
	Tables       []string            // all tables

	Driver SqlDriver // impl SqlDriver instance

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

		if len(tableName) < 0 {
			continue
		}
		tableName = CamelCase(tableName, prefix, true)

		columns := convert.Driver.ReadTablesColumns(tableRealName)
		content := convert.build(tableName, tableRealName, prefix, columns)
		convert.writeModel(tableRealName, content) //写文件
	}
}

//build content with table info
func (convert *Convert) build(tableName, tableRealName, prefix string, columns []Column) (content string) {
	depth := 1
	format := GetFormat(convert.Style)

	content += "package " + convert.PackageName + "\n\n" //写包名
	content += format.AutoImport(tableName)
	content += "type " + tableName + " struct {\n"

	primaryKey := ""
	var primaryColumns Column
	for _, v := range columns {
		var comment string
		if v.ColumnComment != "" {
			comment = fmt.Sprintf(" // %s", v.ColumnComment)
		}
		content += fmt.Sprintf("%s%s %s %s%s\n",
			Tab(depth), v.GetGoColumn(prefix, true), v.GetGoType(), v.GetTag(format), comment)

		if v.IsPrimaryKey() {
			primaryKey = v.ColumnName
			primaryColumns = v
		}

	}

	content += Tab(depth-1) + "}\n\n"

	if primaryKey != "" {
		content += fmt.Sprintf("//get real primary key name \nfunc (%s *%s) %s() string {\n",
			LcFirst(tableName), tableName, "GetKey")
		content += fmt.Sprintf("%sreturn \"%s\"\n",
			Tab(depth), primaryKey)
		content += "}\n\n\n"

		content += fmt.Sprintf("//get primary key in model\nfunc (%s *%s) %s() %s {\n",
			LcFirst(tableName), tableName, "GetKeyProperty", primaryColumns.GetGoType())
		content += fmt.Sprintf("%sreturn %s.%s\n",
			Tab(depth), LcFirst(tableName), CamelCase(primaryKey, prefix, true))
		content += "}\n\n\n"

		content += fmt.Sprintf("//set primary key \nfunc (%s *%s) %s(id %s) {\n",
			LcFirst(tableName), tableName, "SetKeyProperty", primaryColumns.GetGoType())
		content += fmt.Sprintf("%s %s.%s = id\n",
			Tab(depth), LcFirst(tableName), CamelCase(primaryKey, prefix, true))
		content += "}\n\n\n"
	}
	content += fmt.Sprintf("//get real table name\nfunc (%s *%s) %s() string {\n",
		LcFirst(tableName), tableName, "GetTableName")
	content += fmt.Sprintf("%sreturn \"%s\"\n",
		Tab(depth), tableRealName)
	content += "}\n\n\n"
	return content
}

//write file
func (convert *Convert) writeModel(name, content string) {
	log.Printf("write model file %s start\n", name)
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
	log.Printf("write model file %s success\n", name)
}

func (convert *Convert) SetStyle(name string) {
	convert.Style = name
}

func (convert *Convert) GetStyle() string {
	if convert.Style == "" {
		return "default"
	}

	return convert.Style
}

func GetDriver(dir, driver, dsn, style, packageName string) *Convert {
	convert := &Convert{}
	convert.SetPackageName(packageName)
	convert.SetModelPath(dir)

	switch driver {
	case "mysql":
		convert.Driver = &MysqlToGo{}
		convert.Driver.SetDsn(dsn)
		convert.SetStyle(style)
	default:
		panic(fmt.Sprintf("do not support this driver: %v\n", driver))
	}

	return convert
}
