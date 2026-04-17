package aesthetic

import (
	"crypto/md5"
	"encoding/hex"
)

func hashPasswordMD5(password string) string {
	sum := md5.Sum([]byte(password))
	return hex.EncodeToString(sum[:])
}
