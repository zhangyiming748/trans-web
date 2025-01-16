package util

import (
	"strconv"
	"strings"
)

func GenerateUrl(start, offset int, prefix string) {
	var links []string
	for i := 0; i <= offset; i++ {
		// 生成当前的数字
		num := start + i
		// 将数字转换为字符串并写入字符串构建器
		link := strings.Join([]string{prefix, strconv.Itoa(num)}, "/")
		links = append(links, link)
	}
	WriteByLine("post.link", links)
}
