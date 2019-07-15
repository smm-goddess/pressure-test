package text

import (
	"bytes"
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	return generateRandomString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", length)
}

func GenerateRandomStringUpper(length int) string {
	return generateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", length)
}

func GenerateRandomStringLower(length int) string {
	return generateRandomString("abcdefghijklmnopqrstuvwxyz", length)
}

func GenerateRandomStringWithNumbers(length int) string {
	return generateRandomString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", length)
}

func GenerateRandomStringUpperWithNumbers(length int) string {
	return generateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", length)
}

func GenerateRandomStringLowerWithNumbers(length int) string {
	return generateRandomString("abcdefghijklmnopqrstuvwxyz0123456789", length)
}

func generateRandomString(seed string, length int) string {
	source := []byte(seed)
	var buffer bytes.Buffer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		buffer.WriteByte(source[r.Intn(len(source))])
	}
	return buffer.String()
}
