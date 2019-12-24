package generator

import (
	"fmt"
	"strings"
)

type Column struct {
	ColumnName      string
	Type            string
	Nullable        string
	TableName       string
	ColumnComment   string
	Tag             string
	MaxLength       int64
	NumberPrecision int64
	ColumnType      string
	ColumnKey       string
}

func (c Column) GetTag(format Format) string {

	propertyString := c.getProperty(format)
	if propertyString != "" {
		propertyString = strings.TrimRight(propertyString, ";")
	}

	json := c.Tag
	if !format.JsonUseCamel {
		json = CaseCamel(json)
	}

	value := fmt.Sprintf(format.GetTabFormat(), c.Tag, propertyString, json)
	if value != "" {
		if propertyString == "" {
			index := strings.Index(value, ";")
			if index > -1 {
				value = value[0:index] + value[index+1:]
			}
		}
	}

	return value
}

func (c Column) GetGoType() string {
	v, ok := TypeMappingMysqlToGo[c.Type]
	if ok {
		if strings.Index(v, "int") > -1 {
			if strings.Index(c.ColumnType, "unsigned") > -1 {
				v = "u" + v
			}
		}
		return v
	}
	return ""
}

func (c Column) GetMysqlType() string {
	return c.Type
}

func (c Column) GetGoColumn(prefix string, ucFirst bool) string {
	return CamelCase(c.ColumnName, prefix, ucFirst)
}

func (c Column) getProperty(format Format) string {
	if &format.PropertyFormat == nil {
		return ""
	}

	pf := format.GetPropertyFormat()
	value := ""
	var size int64

	if c.MaxLength > 0 {
		size = c.MaxLength
	} else if c.NumberPrecision > 0 {
		size = c.NumberPrecision
	}

	szFormat := pf.GetSizeFormat()
	if size > 0 {
		if szFormat != "" {
			value += fmt.Sprintf(szFormat, size)
			value += ";"
		}
	}

	tpFormat := pf.GetTypeFormat()
	if tpFormat != "" {
		//only support time type
		if strings.Index(strings.ToLower(c.ColumnType), "time") > -1 {
			value += fmt.Sprintf(tpFormat, c.ColumnType)
			value += ";"
		}
	}

	return value
}

func (c Column) IsPrimaryKey() bool {
	return c.ColumnKey == "PRI"
}
