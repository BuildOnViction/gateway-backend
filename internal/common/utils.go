package common

import "encoding/hex"

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
