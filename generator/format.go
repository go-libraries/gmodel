package generator

import "fmt"

type Format struct {
	Framework      string
	TabFormat      string // format must use 3 %s in it, first column name, second property  third json name
	AutoInfo       string
	PropertyFormat PropertyFormat // like size s
}
type PropertyFormat struct {
	Size  string
	Type  string
	Index string
}

//`gorm:"column:beast_id"`
var BeeFormat Format
var DefaultFormat Format
var GormFormat Format

func init() {
	BeeFormat = Format{
		Framework: "bee",
		TabFormat: "`orm:\"column(%s);%s\" json:\"%s\"`",
		PropertyFormat: PropertyFormat{
			Size:  "size(%d)",
			Type:  "type(%s)",
			Index: "%s",
		},
		AutoInfo: "\nimport \"github.com/astaxie/beego/orm\"\n\nfunc init(){\n\torm.RegisterModel(new(%s))\n}\n\n",
	}
	DefaultFormat = Format{
		Framework: "default",
		TabFormat: "`orm:\"%s;%s\" json:\"%s\"`",
	}
	GormFormat = Format{
		Framework: "gorm",
		PropertyFormat: PropertyFormat{
			Size:  "size:%d",
			Type:  "type:%s",
			Index: "",
		},
		TabFormat: "`gorm:\"column:%s;%s\" json:\"%s\"`",
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

func (format Format) AutoImport(modelName string) string {
	if format.AutoInfo == "" {
		return ""
	}
	return fmt.Sprintf(format.AutoInfo, modelName)
}

func (format Format) GetTabFormat() string {
	return format.TabFormat
}

func (format Format) GetPropertyFormat() PropertyFormat {
	return format.PropertyFormat
}

func (pf PropertyFormat) GetSizeFormat() string {
	return pf.Size
}

func (pf PropertyFormat) GetIndexFormat() string {
	return pf.Index
}

func (pf PropertyFormat) GetTypeFormat() string {
	return pf.Type
}
