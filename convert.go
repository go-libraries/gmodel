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
	// 组装struct
	for _, tableRealName := range convert.Tables {
		// 去除前缀
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
			// 字符长度大于1时
			tableName = strings.ToUpper(tableName[0:1]) + tableName[1:]
		}
		depth := 1
		var structContent string
		structContent += "package " + convert.PackageName + "\n\n"
		structContent += "type " + tableName + " struct {\n"
		columns, ok := convert.TableColumn[tableRealName]
		for _, v := range columns {
			//structContent += tab(depth) + v.ColumnName + " " + v.Type + " " + v.Json + "\n"
			// 字段注释
			var comment string
			if v.ColumnComment != "" {
				comment = fmt.Sprintf(" // %s", v.ColumnComment)
			}
			structContent += fmt.Sprintf("%s%s %s %s%s\n",
				Tab(depth), v.GetGoColumn(prefix, true), v.GetGoType(), v.GetTag("orm"), comment)
		}
		structContent += Tab(depth-1) + "}\n\n"

		// 添加 method 获取真实表名
		structContent += fmt.Sprintf("func (%s *%s) %s() string {\n",
			LcFirst(tableName), tableName, "GetTableName")
		structContent += fmt.Sprintf("%sreturn \"%s\"\n",
			Tab(depth), tableRealName)
		structContent += "}\n\n"
		convert.writeModel(tableRealName, structContent)
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
