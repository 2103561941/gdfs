package util

import (
	"crypto/md5"
	"fmt"
)

func EncodeMD5(str string) string {
	m5 := md5.New()
	m5.Write([]byte(str))
	st := m5.Sum(nil)
	md5str := fmt.Sprintf("%x", st) //将[]byte转成16进制
	return md5str
}
