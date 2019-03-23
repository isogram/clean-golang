package utils

import (
	"crypto/md5"
	"fmt"
)

// StringMD5 generating MD5 string
func StringMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
