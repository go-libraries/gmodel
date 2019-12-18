

/*
 * Copyright (c) 2019 Mars Lee. All rights reserved.
 */

package model

import "strings"

// aaa_bbb => AaaBbb
// if !ucFirst
// aaa_bbb => aaabbb
func CamelCase(str,prefix string, ucFirst bool) string {
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