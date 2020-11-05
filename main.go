package main

import (
	"flag"
	"fmt"
	"github.com/go-libraries/genModels"
	"os"
	"strings"
)

var (
	modelPath   string
	driver      string
	dsn         string
	ignoreTable string
	style       string
	packageName string
	camel       string
	help        bool
	h           bool
)

func init() {
	currentPath, _ := os.Getwd()
	flag.StringVar(&modelPath, "dir", currentPath, "a dir name save model file path, default is current path")
	flag.StringVar(&driver, "driver", "mysql", "database driver,like `mysql` `mariadb`, default is mysql")
	flag.StringVar(&dsn, "dsn", "", "connection info names dsn")
	flag.StringVar(&ignoreTable, "ig_tables", "", "ignore table names, like [tableA,tableB]")
	flag.StringVar(&style, "style", "default", "use orm style like `bee` `gorm`, default `default`")
	flag.StringVar(&camel, "camel", "0", "json use camel 0 false 1 true")
	flag.StringVar(&packageName, "package", "", "help")
	flag.BoolVar(&help, "help", false, "this help")
	flag.BoolVar(&h, "h", false, "this help")
}

func main() {
	flag.Parse()
	if h || help {
		flag.Usage()
	}
	//dsn = "ceshiceshi:mgtj123456@tcp(rm-ly29w98y58jjy4s7x.mysql.rds.aliyuncs.com:3306)/mgtj_app_content?charset=utf8mb4"
	//driver = "mysql"
	//style = "grom"
	//packageName = "content"
	//modelPath = "d:\\work\\contents"
	if dsn == "" {
		flag.Usage()
	}

	if modelPath != "" {
		_,e := os.Stat(modelPath)
		if e != nil {
			_ = os.Mkdir(modelPath, os.ModePerm)
		}
	}
	flag.Usage = usage
	genModels.GormFormat.JsonUseCamel = false
	if camel == "0" {
		genModels.GormFormat.JsonUseCamel = true
	}
	fmt.Println(genModels.GormFormat,camel, style)
	driver := genModels.GetDriver(modelPath, driver, dsn, style, packageName)
	if ignoreTable != "" {
		ignoreTables := strings.Split(ignoreTable, ",")
		driver.SetIgnoreTables(ignoreTables...)
	}
	driver.Run()
}

func usage() {
	fmt.Println("Usage: model [-dir dirname] [-driver mysql] [-dsn username:password@tcp(host:port)/database]")
	flag.PrintDefaults()
}
