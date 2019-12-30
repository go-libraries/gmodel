package generator

import (
	"strings"
)

type Format struct {
	Framework      string
	TabFormat      string // format must use 3 %s in it, first column name, second property  third json name
	AutoInfo       string
	PropertyFormat PropertyFormat // like size s
	JsonUseCamel   bool
}
type PropertyFormat struct {
	Size  string
	Type  string
	Index string
	Default string
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
			Default: "",
		},
		AutoInfo: "\nimport \"github.com/astaxie/beego/orm\"\n\nfunc init(){\n\torm.RegisterModel(new({{modelName}}))\n}\n\n",
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
			Default: "default:%s",
		},
		TabFormat: "`gorm:\"column:%s;%s\" json:\"%s\"`",
		//AutoInfo:  "\nimport (\n\t\"fmt\"\n)\n\n",
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
	return strings.Replace(format.AutoInfo, "{{modelName}}", modelName, -1)
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

func (pf PropertyFormat) GetDefaultFormat() string {
	return pf.Default
}

var GormTpl = `
func ({{entry}} *{{object}}) GetById(id string) {
	Orm.Model({{entry}}).First({{entry}}, {{entry}}.GetKey() + " = '"+id+"'")
}


func ({{entry}} *{{object}}) GetOne(condition string) (err []error) {
	err = Orm.Model({{entry}}).First(userInfo, condition).GetErrors()
	if len(err) > 0 {
		return err
	}
	return
}

func ({{entry}} *{{object}}) GetList(page,limit int64, condition interface{}) (list []*{{object}}) {
	err := Orm.Model({{entry}}).Limit(limit).Offset((page-1) * limit).Find(list, condition).GetErrors()
	if err != nil {
		return nil
	}
	return
}

func ({{entry}} *{{object}}) Create() []error {
	return Orm.Model({{entry}}).Create({{entry}}).GetErrors()
}

func ({{entry}} *{{object}}) Update(info UserInfo) []error  {
	return Orm.Model({{entry}}).UpdateColumns(info).GetErrors()
}

func ({{entry}} *{{object}}) Delete()  {
	Orm.Model({{entry}}).Delete({{entry}})
}
`
var GormInit = `
package {{package}}

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Orm *gorm.DB

func init() {
	db, err := gorm.Open("mysql", "{{dns}}")
	if err != nil {
		panic("连接数据库失败")
	}
	Orm = db
}
`
func (format Format) GetFuncTemplate(t string) string  {
	switch t {
	case "gorm":
		return GormTpl
	default:
		return ""
	}
}

func (format Format) GetInitTemplate(t string) string  {
	switch t {
	case "gorm":
		return GormInit
	default:
		return ""
	}
}