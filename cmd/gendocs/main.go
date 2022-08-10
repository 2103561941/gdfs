// 自动生成 cobra 的命令行文档文件

package main

import (
	"gdfs/internal/client/cmd"
	"log"
)

func main() {
	filepath := "api/cmd/"
	log.Println(cmd.GenDocs(filepath))
}