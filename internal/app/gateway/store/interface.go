package database

import . "github.com/anhntbk08/gateway/internal/app/gateway/store/entity"

// type Databaser interface {
// 	RequestLoginToken(address string) (interface{}, error)
// 	Login(authService.Token) string
// }

type UserDao interface {
	// IssueToken(address string) error
	// SaveLoginSession(authService.Token) (authService.Token, error)
	Upsert(user *User) error // update or insert new user
	IsAuthen(token string) (*User, error)
	IsExist(address string) (*User, error)
}
