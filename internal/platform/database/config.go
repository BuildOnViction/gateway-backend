package database

import (
	"emperror.dev/errors"
)

// Config holds information necessary for connecting to a database.
type Config struct {
	Uri    string
	DbName string
	Params map[string]string
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	if c.Uri == "" {
		return errors.New("database uri is required")
	}

	return nil
}
