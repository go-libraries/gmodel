package model

import "fmt"

type column struct {
	ColumnName    string
	Type          string
	Nullable      string
	TableName     string
	ColumnComment string
	Tag           string
}

func (c column) GetTag(tagKey string) string {
	return fmt.Sprintf("`%s:\"%s\" json:\"%s\"`", tagKey, c.Tag, c.Tag)
}

func (c column) GetGoType() string {
	v,ok := TypeMappingMysqlToGo[c.Type]; if ok {
		return v
	}
	return ""
}

func (c column) GetMysqlType() string  {
	return c.Type
}

func (c column) GetGoColumn(prefix string, ucFirst bool) string  {
	return CamelCase(c.ColumnName, prefix, ucFirst)
}
