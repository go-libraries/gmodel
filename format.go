package model

import "fmt"

type Format struct {
	Framework string
	TabFormat string  // format must use 2 %s in it, first column name, second json name
	AutoInfo string
}
//`gorm:"column:beast_id"`
var BeeFormat Format
var DefaultFormat Format
var GormFormat Format
func init()  {
	BeeFormat = Format{
		Framework: "bee",
		TabFormat: "`orm:\"column(%s)\" json:\"%s\"`",
		AutoInfo: "\nimport \"github.com/astaxie/beego/orm\"\n\nfunc init(){\n\torm.RegisterModel(new(%s))\n}\n",
	}
	DefaultFormat = Format{
		Framework: "default",
		TabFormat: "`orm:\"%s\" json:\"%s\"`",
	}
	GormFormat = Format{
		Framework: "gorm",
		TabFormat: "`gorm:\"column:%s\" json:\"%s\"`",
		AutoInfo:  "",
	}
}

func GetFormat(framework string) Format {
	switch framework {
	case "bee":
		return BeeFormat
	case "gorm":
		return GormFormat
	default:
		return DefaultFormat
	}
}

func (format Format) AutoImport(modelName string) string  {
	if format.AutoInfo == "" {
		return ""
	}
	return fmt.Sprintf(format.AutoInfo, modelName)
}

func (format Format) GetTabFormat()  string {
	return format.TabFormat
}