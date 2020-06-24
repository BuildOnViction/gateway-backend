package auth

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"emperror.dev/errors"
	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/jwt"
	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	. "github.com/anhntbk08/gateway/internal/common"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"github.com/tomochain/tomochain/crypto"
	"golang.org/x/crypto/sha3"
	// . "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
)

// +kit:endpoint:errorStrategy=auth

type Service interface {
	RequestToken(ctx context.Context, request RqTokenData) (token Token, err error)
	Login(ctx context.Context, request Token) (accessToken string, err error)
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
	db      *Mongo
	authKey string
}

func NewService(db *Mongo, authKey string) Service {
	return &service{
		db:      db,
		authKey: authKey,
	}
}

func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func verifySig(from, sigHex, msg string) (bool, error) {
	byteMsg := []byte(msg)
	fromAddr := common.HexToAddress(from)
	sig := hexutil.MustDecode(sigHex)
	if sig[64] != 27 && sig[64] != 28 {
		return false, errors.New("Wrong format signature")
	}
	sig[64] -= 27
	pubKey, err := crypto.SigToPub(signHash(byteMsg), sig)
	if err != nil {
		return false, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return strings.ToLower(fromAddr.Hex()) == strings.ToLower(recoveredAddr.Hex()), nil
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

func (s service) Login(ctx context.Context, request Token) (accessToken string, err error) {
	if !IsValidAddress(request.Address) {
		return "", errors.WithStack(ValidationError{Violates: map[string][]string{
			"address": {
				"AUTH.REQUEST_TOKEN.INVALID_ADDRESS",
				"Invalid address " + err.Error(),
			},
		}})
	}

	if _, err := s.db.SessionDao.IsValidToken(request.Address, request.Token); err != nil {
		return "", errors.WithStack(ValidationError{Violates: map[string][]string{
			"token": {
				"AUTH.LOGIN.INVALID_TOKEN",
				"Invalid token " + err.Error(),
			},
		}})
	}

	if valid, err := verifySig(request.Address, request.Signature, request.Token); !valid {
		return "", errors.WithStack(ValidationError{Violates: map[string][]string{
			"signature": {
				"AUTH.LOGIN.INVALID_SIGNATURE",
				"Invalid signature " + err.Error(),
			},
		}})
	}
	logintoken, err := GenerateToken([]byte(s.authKey), request.Address)

	if err != nil {
		return "", err
	}

	s.db.SessionDao.Used(request.Token)
	err = s.db.UserDao.Upsert(&entity.User{
		Address: request.Address,
		Session: entity.AuthenSession{
			ExpiredAt: time.Now().Add(time.Hour * 24),
			Token:     logintoken,
		},
		UpdatedAt: time.Now(),
	})

	return logintoken, err
}
