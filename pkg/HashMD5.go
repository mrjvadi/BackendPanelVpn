package pkg

import (
	"crypto/md5"
	"encoding/hex"
)

func HashMD5(password string) string {
	// محاسبهٔ چک‌سام 16 بایتی
	sum := md5.Sum([]byte(password))
	// تبدیل بایت‌ها به رشتهٔ هگزادسیمال
	return hex.EncodeToString(sum[:])
}
