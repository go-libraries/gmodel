/*
 * Copyright (c) 2019 Mars Lee. All rights reserved.
 */

package generator

import (
	"strconv"
	"strings"
)

// aaa_bbb => AaaBbb
// if !ucFirst
// aaa_bbb => aaabbb
func CamelCase(str, prefix string, ucFirst bool) string {
	// 是否有表前缀, 设置了就先去除表前缀
	if prefix != "" {
		str = strings.Replace(str, prefix, "", 1)
	}
	var text string
	//for _, p := range strings.Split(name, "_") {
	for _, p := range strings.Split(str, "_") {
		// 字段首字母大写的同时, 是否要把其他字母转换为小写
		switch len(p) {
		case 0:
		case 1:
			text += strings.ToUpper(p[0:1])
		default:
			// 字符长度大于1时
			if ucFirst == true {
				text += UcFirst(p[0:1]) + LcFirst(p[1:])
			} else {
				text += UcFirst(p[0:1]) + p[1:]
			}
		}
	}
	return text
}

func CaseCamel(str string) string {

	var text string

	for _, p := range str[:] {
		if p > 64 && p < 91 {
			p = p + 32
			text += "_"
		}
		text += string(p)
	}

	return text
}

//format
func Tab(depth int) string {
	return strings.Repeat("\t", depth)
}

func UcFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func LcFirst(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func Byte2Int64(data []byte) int64 {
	var str string
	var ret int64 = 0
	for i := 0; i < len(data); i++ {
		str += string(data[i])
	}
	ret, _ = strconv.ParseInt(str, 10, 64)
	return ret
}
