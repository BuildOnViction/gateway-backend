package common

import (
	"encoding/hex"
	"regexp"
)

func IsValidMongoID(id string) bool {
	idByte, err := hex.DecodeString(id)
	if err != nil {
		return false
	}

	if len(idByte) != 12 {
		return false
	}

	return true
}

func IsAddress(address string) bool {
	var rex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

	return len(rex.FindAllString(address, -1)) > 0
}
