package golang

import (
	"encoding/base32"
	"encoding/binary"
	"math/rand"
	"time"
	"unsafe"
)

const RAND_SIZE = 32

var randArray [RAND_SIZE]byte

// GenerateRandom - выдает случайную строку
func GenerateRandom() string {
	return QuickRandom(32)
}

func SimplestRandom() string {
	for i := 0; i < RAND_SIZE/8; i++ {
		binary.BigEndian.PutUint64(randArray[i*8:i*8+8], rand.Uint64())
	}
	return base32.StdEncoding.EncodeToString(randArray[:])
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func QuickRandom(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
