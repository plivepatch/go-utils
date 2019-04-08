package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Md5(raw string) string {
	h := md5.Sum([]byte(raw))
	return hex.EncodeToString(h[:])
}

func CalcMd5(data string, protocol string, salt string) string {
	str := data + salt + protocol
	hash := md5.Sum([]byte(str))
	md5now := fmt.Sprintf("%X", hash)
	return md5now
}

func VerifyMd5(data string, protocol string, salt string, md5str string) (bool, string) {
	str := data + salt + protocol
	hash := md5.Sum([]byte(str))
	md5now := fmt.Sprintf("%X", hash)
	if md5str == md5now {
		return true, md5str
	} else {
		return false, md5str
	}
}
