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

type ChainConfig struct {
	RPC string
	WS  string

	SignMode            string
	Confirmations       int64
	StartBlock          uint64
	BlockTime           time.Duration
	IntervalRunningTime time.Duration
	Unit                int64
	NeededConfirms      int64
	GasLimit            uint64
	MinAcceptedValue    int64
	MinGasPrice         uint64
}
