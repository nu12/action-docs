package helper

import (
	"crypto/md5"
	"fmt"
)

func Hash(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}
