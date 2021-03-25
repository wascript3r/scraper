package sha256

import (
	"crypto/hmac"
	"crypto/sha256"
)

const (
	HashLen = 32
)

func Compute(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func ComputeHMAC(message, secret []byte) ([]byte, error) {
	h := hmac.New(sha256.New, secret)

	_, err := h.Write(message)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
