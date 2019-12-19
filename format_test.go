package model

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestBeeFormat(t *testing.T) {
	tab := BeeFormat.GetTabFormat()

	tab = fmt.Sprintf(tab, "columnName", "columnName")

	assert.Equal(t, tab, "`orm:\"column(columnName)\" json:\"columnName\"`", "\nerror!!! not equal")
}