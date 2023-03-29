package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func CIntn(ln, n int) []int {
	c := make([]int, ln)
	for j := range c {
		for i := 0; i < 10; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
			if err != nil {
				continue
			}
			c[j] = int(num.Int64())
		}
	}
	return c
}

func String(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	read, err := rand.Read(bytes)
	// Note that err == nil iff we read len(b) bytes.
	if err != nil {
		return "", err
	}
	if read != n {
		return "", fmt.Errorf("failed to read %d bytes from random, bytes read %d", n, read)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func InvitationCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i, j := range CIntn(n, len(letters)) {
		b[i] = letters[j]
	}
	return string(b)
}

func NumString(n int) string {
	b := strings.Builder{}
	for _, v := range CIntn(n, 10) {
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}
