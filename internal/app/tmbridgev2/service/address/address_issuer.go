package address

import (
	"errors"
	"strings"
	"sync"

	// . "github.com/anhntbk08/gateway/internal/app/tmbridgev2/jwt"
	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomochain/crypto"
	"github.com/wemeetagain/go-hdwallet"
)

type AddressIssuer struct {
	db          *Mongo
	xPublicKeys map[string]string
	mux         sync.Mutex
}

func NewIssuer(db *Mongo, xPublicKeys map[string]string) *AddressIssuer {
	return &AddressIssuer{db: db, xPublicKeys: xPublicKeys}
}

func (issuer *AddressIssuer) getAddress(projectID bson.ObjectId, coin string) (string, uint64, error) {
	issuer.mux.Lock()
	defer issuer.mux.Unlock()

	coin = strings.ToLower(coin)

	if issuer.xPublicKeys[coin] == "" {
		return "", 0, errors.New("Not support " + coin + " yet")
	}

	var addressDao entity.Address
	err := issuer.db.AddressDao.GetSortOne(bson.M{"coin_type": coin, "project": projectID}, []string{"-_id"}, &addressDao)

	if err != nil {
		return "", 0, err
	}

	coinAccountIndex := addressDao.AccountIndex + 1
	var address string

	if coin == "btc" {
		address, err = issuer.getBTCAddress(coinAccountIndex, coin)
		if err != nil {
			return "", 0, err
		}
	} else if coin == "eth" {
		address, err = issuer.getEthAddress(coinAccountIndex, coin)
		if err != nil {
			return "", 0, err
		}
	} else if coin == "usdt" {
		address, err = issuer.getEthAddress(coinAccountIndex, coin)
		if err != nil {
			return "", 0, err
		}
	} else {
		return "", 0, errors.New("Not support " + coin + " yet")
	}

	return address, coinAccountIndex, nil
}

func (issuer *AddressIssuer) getBTCAddress(coinAccountIndex uint64, coin string) (string, error) {
	childPubKey, err := hdwallet.StringChild(
		issuer.xPublicKeys[coin], uint32(coinAccountIndex), // we're always sure that total account < 2^32 bit
	)

	if err != nil {
		return "", err
	}

	return hdwallet.StringAddress(childPubKey)
}

func (issuer *AddressIssuer) getEthAddress(coinAccountIndex uint64, coin string) (string, error) {
	childPubKey, err := hdwallet.StringChild(
		issuer.xPublicKeys[coin], uint32(coinAccountIndex), // we're always sure that total account < 2^32 bit
	)
	if err != nil {
		return "", err
	}
	w, err := hdwallet.StringWallet(childPubKey)
	if err != nil {
		return "", err
	}
	ecdsaPubKey, err := crypto.DecompressPubkey([]byte(w.Key))

	if err != nil {
		return "", err
	}

	address := crypto.PubkeyToAddress(*ecdsaPubKey).Hex()

	return address, nil
}

func (issuer *AddressIssuer) getTomoAddress(coinAccountIndex uint64, coin string) (string, error) {
	return issuer.getEthAddress(coinAccountIndex, coin)
}
