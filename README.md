# model-gen

每次写go的orm的时候，总要去来回粘贴数据表，然后orm的东西又不对应，为此，懒人李只能造个轮子，为了提高效率(ps:就是想偷懒、摸鱼)。

一个Go专用的模型生成器就这样诞生啦

# 原理

通过获取mysql表结构，进行model文件生成。

# 迭代

目前支持mysql，未来预计支持mariadb和pgsql(sql server还未考量)

# 快速入门

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


