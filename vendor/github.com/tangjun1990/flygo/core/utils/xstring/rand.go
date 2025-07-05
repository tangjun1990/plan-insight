package xstring

import (
	"encoding/hex"
	"math/rand"
	"time"
)

var _randSource = rand.New(rand.NewSource(time.Now().Unix()))

// RandByte 随机长度的byte
func RandByte(byteLen int) []byte {
	id := make([]byte, byteLen)
	_, _ = _randSource.Read(id)
	return id
}

// RandByteStr 随机长度的byte 字符串
func RandByteStr(byteLen int) string {
	return hex.EncodeToString(RandByte(byteLen))
}
