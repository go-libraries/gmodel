# model-gen

每次写go的orm的时候，总要去来回粘贴数据表，然后orm的东西又不对应，为此，懒人李只能造个轮子，为了提高效率(ps:就是想偷懒、摸鱼)。

一个Go专用的模型生成器就这样诞生啦

# 原理

通过获取数据库表结构，进行model文件生成。

# 迭代

目前支持mysql，未来预计支持mariadb和pgsql(sql server还未考量)

2019-12-19 

1. 支持多种风格文件生成(目前支持 bee、gorm、默认格式)
2. 新增无符号整型和标签内容（size、type）

2019-12-23

1. 增加主键支持

2019-12-24

1. 增加gorm curd方法
2. 增加orm自动初始化模板并支持orm


# 快速入门


    go get -u github.com/go-libraries/gmodel
    
    gmodel -dsn="root:sa@tcp(localhost:3306)/blog" -dir=/tmp -style=bee -package=model
    
    cat /tmp/cate.go
   
文件内容如下：   
```go
package model

import "github.com/astaxie/beego/orm"

func init(){
	orm.RegisterModel(new(Cate))
}

type Cate struct {
	Id         uint   `orm:"column(id);size(10);type(int(11) unsigned);" json:"id"`
	Name       string `orm:"column(name);size(50);type(varchar(50));" json:"name"`
	CreateTime string `orm:"column(create_time);type(timestamp);" json:"create_time"`
	UpdateTime string `orm:"column(update_time);type(timestamp);" json:"update_time"`
}

func (cate *Cate) GetTableName() string {
	return "cate"
}
```
s
help:
    
    Usage of gmodel:
      -dir string
        	a dir name save model file path, default is current path (default "/Users/limars/Go/bin")
      -driver mysql
        	database driver,like mysql `mariadb`, default is mysql (default "mysql")
      -dsn string
        	connection info names dsn
      -h	this help
      -help
        	this help
      -ig_tables string
        	ignore table names
      -package string
        	help
      -style bee
        	use orm style like bee `gorm`, default `default` (default "default")

# 代码应用

代码：
```
	    mysqlHost := "127.0.0.1"
    	mysqlPort := "3306"
    	mysqlUser := "root"
    	mysqlPassword := "sa"
    	mysqlDbname := "blog"
    
    	dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDbname + "?charset=utf8mb4"
    
    	Mysql := GetMysqlToGo()
    	Mysql.Driver.SetDsn(dsn)
    	Mysql.SetModelPath("/tmp")
    	Mysql.SetIgnoreTables("cate")
    	Mysql.SetPackageName("models")
    	Mysql.Run()
```

执行结果:

    ll /tmp
    
```
total *
-rw-r--r--  1 limars  wheel  297 12 18 17:59 cate.go
-rw-r--r--  1 limars  wheel  597 12 18 17:59 comment.go
-rw-r--r--  1 limars  wheel  826 12 18 17:59 content.go
......
```

    cat cate.go
```
package models

type Cate struct {
	Id         int    `orm:"id" json:"id"`
	Name       string `orm:"name" json:"name"`
	CreateTime string `orm:"create_time" json:"create_time"`
	UpdateTime string `orm:"update_time" json:"update_time"`
}

func (cate *Cate) GetTableName() string {
	return "cate"
}
```


附带上sql:

```
CREATE DATABASE `blog`;

USE `blog`;

CREATE TABLE `cate` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `comment` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `blog_id` int(11) unsigned NOT NULL,
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0',
  `ip` varchar(32) NOT NULL DEFAULT '',
  `email` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(50) NOT NULL DEFAULT '',
  `content` tinytext NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `content` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `cate_id` int(11) unsigned NOT NULL COMMENT '分类id',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `description` tinytext NOT NULL COMMENT '简介',
  `content` text NOT NULL COMMENT '正文',
  `keyword` varchar(255) NOT NULL DEFAULT '' COMMENT 'seo关键字',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `is_original` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 原创  2 转载',
  `ext` text NOT NULL COMMENT '扩展字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```


