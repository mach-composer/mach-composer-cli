package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func ComputeHash(t interface{}) (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}
