package common

type JWT struct {
	ExpiredTime uint64 // in minutes
	Key         string
}
