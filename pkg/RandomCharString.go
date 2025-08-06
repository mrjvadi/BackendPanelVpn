package pkg

import (
	"crypto/rand"
	"math/big"
)

// Alphabet مجموعهٔ کاراکترهایی است که می‌خواهیم از آن‌ها انتخاب کنیم.
// می‌توانید هر کاراکتر دیگری هم اضافه کنید.
const Alphabet = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

// RandomString تولید یک رشتهٔ تصادفی به طول n از Alphabet.
func RandomString(n int) string {
	result := make([]byte, n)
	alphaLen := big.NewInt(int64(len(Alphabet)))
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, alphaLen)
		if err != nil {
			return ""
		}
		result[i] = Alphabet[idx.Int64()]
	}
	return string(result)
}
