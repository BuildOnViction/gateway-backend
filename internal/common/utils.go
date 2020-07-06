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

func Difference(slice1 []string, slice2 []string) (inSlice1 []string, inSlice2 []string) {
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			inSlice1 = append(inSlice1, s1)
		}
	}

	for _, s2 := range slice2 {
		found := false
		for _, s1 := range slice1 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			inSlice2 = append(inSlice2, s2)
		}
	}

	return inSlice1, inSlice2
}
