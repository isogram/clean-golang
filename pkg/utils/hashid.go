package utils

import (
	"os"

	hashids "github.com/speps/go-hashids"
)

const (
	// HasIDMinLength minimum hasID length
	HasIDMinLength = 3
)

// HashedToInt64 to convert HashedID to integer64
func HashedToInt64(hashed string) int64 {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("APP_KEY")
	hd.MinLength = HasIDMinLength
	h, _ := hashids.NewWithData(hd)
	decoded := h.DecodeInt64(hashed)

	return decoded[0]
}

// Int64ToHashed hash integer to string with hashid
func Int64ToHashed(i int64) string {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("APP_KEY")
	hd.MinLength = HasIDMinLength
	h, _ := hashids.NewWithData(hd)
	encoded, _ := h.EncodeInt64([]int64{i})
	return encoded
}
