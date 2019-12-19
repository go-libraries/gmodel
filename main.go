package main

import (
	"flag"
	"fmt"
	"github.com/go-libraries/gmodel/generator"
	"os"
)

var (
	modelPath   string
	driver      string
	dsn         string
	ignoreTable string
	style       string
	packageName string
	help        bool
	h           bool
)

func init() {
	currentPath, _ := os.Getwd()
	flag.StringVar(&modelPath, "dir", currentPath, "a dir name save model file path, default is current path")
	flag.StringVar(&driver, "driver", "mysql", "database driver,like `mysql` `mariadb`, default is mysql")
	flag.StringVar(&dsn, "dsn", "", "connection info names dsn")
	flag.StringVar(&ignoreTable, "ig_tables", "", "ignore table names")
	flag.StringVar(&style, "style", "default", "use orm style like `bee` `gorm`, default `default`")
	flag.StringVar(&packageName, "package", "", "help")
	flag.BoolVar(&help, "help", false, "this help")
	flag.BoolVar(&h, "h", false, "this help")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
	}

	if dsn == "" {
		flag.Usage()
	}
	flag.Usage = usage
	generator.GetDriver(modelPath, driver, dsn, style, packageName).Run()
}

func usage() {
	fmt.Println("Usage: model [-dir dirname] [-driver mysql] [-dsn username:password@tcp(host:port)/database]")
	flag.PrintDefaults()
}
