package auth

import (
	"context"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	"emperror.dev/errors"

	. "github.com/anhntbk08/gateway/internal/app/gateway/store"
	"github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
	. "github.com/anhntbk08/gateway/internal/common"
	"github.com/tomochain/tomochain/crypto"
	"golang.org/x/crypto/sha3"
	// . "github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
)

// +kit:endpoint:errorStrategy=auth

type Service interface {
	RequestToken(ctx context.Context, request RqTokenData) (token Token, err error)
	Login(ctx context.Context, request Token) (success bool, err error)
}

type RqTokenData struct {
	Address string
}

type Token struct {
	ID        string
	Address   string
	Signature string
	Token     string
}

type service struct {
	db *Mongo
}

func NewService(db *Mongo) Service {
	return &service{
		db: db,
	}
}

func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

func IsValidateSignature(address, token, signature string) (bool, error) {
	byteToken, err := hex.DecodeString(token)
	byteSign, err := hex.DecodeString(signature)

	pubkey, err := crypto.Ecrecover(
		byteToken,
		byteSign,
	)

	if err != nil {
		return false, err
	}

	ecdsaPubkey := crypto.ToECDSAPub(pubkey)
	if err != nil {
		return false, err
	}

	expectedAddress := crypto.PubkeyToAddress(
		*ecdsaPubkey,
	)
	return strings.ToLower(expectedAddress.Hex()) == strings.ToLower(address), nil
}

func (s service) RequestToken(ctx context.Context, request RqTokenData) (token Token, err error) {
	if !IsValidAddress(request.Address) {
		return Token{}, errors.WithStack(ValidationError{Violates: map[string][]string{
			"address": {
				"AUTH.INVALID_ADDRESS",
				"Invalid address",
			},
		}})
	}

	hashFunc := sha3.New224()
	sha3 := func(data ...[]byte) []byte {
		h := hashFunc
		for _, v := range data {
			h.Write(v)
		}
		return h.Sum(nil)
	}

	// TODO put to conf
	hashString := time.Now().String() + "_tomochain_bridgegateway"
	issuedToken := hex.EncodeToString(
		crypto.Keccak256(sha3(
			[]byte(hashString),
		)),
	)

	err = s.db.SessionDao.Create(&entity.Session{
		Address:   request.Address,
		Token:     issuedToken,
		ExpiredAt: time.Now().Add(time.Minute * 10),
	})

	return Token{
		Address: request.Address,
		Token:   issuedToken,
	}, err
}

func (s service) Login(ctx context.Context, request Token) (success bool, err error) {
	if !IsValidAddress(request.Address) {
		return false, errors.WithStack(ValidationError{Violates: map[string][]string{
			"address": {
				"AUTH.REQUEST_TOKEN.INVALID_ADDRESS",
				"Invalid address " + err.Error(),
			},
		}})
	}

	if _, err := s.db.SessionDao.IsValidToken(request.Address, request.Token); err != nil {
		return false, errors.WithStack(ValidationError{Violates: map[string][]string{
			"token": {
				"AUTH.LOGIN.INVALID_TOKEN",
				"Invalid token " + err.Error(),
			},
		}})
	}

	if valid, err := IsValidateSignature(request.Address, request.Token, request.Signature); !valid {
		return false, errors.WithStack(ValidationError{Violates: map[string][]string{
			"signature": {
				"AUTH.LOGIN.INVALID_SIGNATURE",
				"Invalid signature " + err.Error(),
			},
		}})
	}
	s.db.SessionDao.Used(request.Token)
	err = s.db.UserDao.Upsert(&entity.User{
		Address: request.Address,
		Session: entity.AuthenSession{
			Signature: request.Signature,
			ExpiredAt: time.Now().Add(time.Hour * 24),
			Token:     request.Token,
		},
		UpdatedAt: time.Now(),
	})

	return true, err
}
