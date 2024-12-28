package crypto

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID 生成一个新的 UUID
func GenerateUUID() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func GenerateUUIDWithoutHyphen() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(newUUID.String(), "-", ""), nil
}
