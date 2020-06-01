package auth

import (
	"context"
)

// +kit:endpoint:errorStrategy=auth

type Service interface {
	RequestToken(ctx context.Context, request RqTokenData) (token Token, err error)
}

type RqTokenData struct {
	Address string
}
type Token struct {
	ID          string
	Address     string
	Signature   string
	IssuedToken string
}

func NewService(idgenerator IDGenerator) Service {
	return &service{
		idgenerator: idgenerator,
	}
}

type service struct {
	idgenerator IDGenerator
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// NotFoundError is returned if an item cannot be found.
type NotFoundError struct {
	ID string
}

// Error implements the error interface.
func (NotFoundError) Error() string {
	return "item not found"
}

// Details returns error details.
func (e NotFoundError) Details() []interface{} {
	return []interface{}{"item_id", e.ID}
}

// NotFound tells a client that this error is related to a resource being not found.
// Can be used to translate the error to eg. status code.
func (NotFoundError) NotFound() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (NotFoundError) ServiceError() bool {
	return true
}

type validationError struct {
	violations map[string][]string
}

// Validation tells a client that this error is related to a resource being invalid.
// Can be used to translate the error to eg. status code.
func (validationError) Validation() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (validationError) ServiceError() bool {
	return true
}

func (s service) RequestToken(ctx context.Context, request RqTokenData) (token Token, err error) {
	return Token{
		Address:     request.Address,
		IssuedToken: "hellonewcommer",
	}, nil
}
