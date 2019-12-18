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

type Convert struct {
	ModelPath    string
	TablePrefix  map[string]string
	TableColumn  map[string][]column
	IgnoreTables []string
	Tables       []string
	PackageName  string
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

func (convert *Convert) Run() {
	for _, tableRealName := range convert.Tables {
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
		depth := 1
		var content string
		content += "package " + convert.PackageName + "\n\n" //写包名
		content += "type " + tableName + " struct {\n"
		columns, ok := convert.TableColumn[tableRealName]
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
		convert.writeModel(tableRealName, content) //写文件
	}
}

func (convert *Convert) writeModel(name, content string) {
	filePath := fmt.Sprintf("%s/%s.go", convert.ModelPath, name)
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Can not write file" + filePath)
		return
	}
	defer f.Close()

	_,err = f.WriteString(content)
	if err != nil {
		log.Println("Can not write file" + filePath)
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	_ = cmd.Run()
}
