package common

import "time"

type JWT struct {
	ExpiredTime uint64 // in minutes
	Key         string
}

type JobqueueConfig struct {
	Broker          string
	DefaultQueue    string
	ResultBackend   string
	ResultsExpireIn time.Duration
	Mongodb         string
	ProjectID       string
	GoogleAuth      string
	JobPrefix       string
}
