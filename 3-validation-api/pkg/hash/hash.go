package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateRandomHash() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	data := fmt.Sprintf("%d%s", timestamp, string(randomBytes))
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
