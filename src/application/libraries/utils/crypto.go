package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	md5Obj := md5.New()
	md5Obj.Write([]byte(str))
	return hex.EncodeToString(md5Obj.Sum(nil))
}
