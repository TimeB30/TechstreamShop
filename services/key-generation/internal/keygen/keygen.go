package keygen

import (
	"os/exec"
	"strconv"
)

type KeyGenerator struct{}

func NewKeyGen() *KeyGenerator {
	return &KeyGenerator{}
}

func (k KeyGenerator) GenerateKey(SoftwareID string, days int64, version string) (string, error) {
	daysStr := strconv.FormatInt(days, 10)
	cmd := exec.Command("python", "generatekey.py", SoftwareID, daysStr)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
