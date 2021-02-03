package golang

import (
	"encoding/base32"
	"encoding/binary"
	"math/rand"
	"time"
)

const RAND_SIZE = 32

var randArray [RAND_SIZE]byte

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandom - выдает случайную строку
func GenerateRandom() string {
	for i := 0; i < RAND_SIZE/8; i++ {
		binary.BigEndian.PutUint64(randArray[i*8:i*8+8], rand.Uint64())
	}
	return base32.StdEncoding.EncodeToString(randArray[:])
}
