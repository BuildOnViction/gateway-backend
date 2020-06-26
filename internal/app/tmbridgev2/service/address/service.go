package address

import (
	"context"

	"emperror.dev/errors"

	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/globalsign/mgo/bson"
)

// +kit:endpoint:errorStrategy=address

type Service interface {
	Issue(ctx context.Context, issueRequest IssueRequest) (entity.Address, error)

	List(ctx context.Context, listRequest ListRequest) ([]entity.Address, error)

	IsIsssueBy(ctx context.Context, address string) (bool, error)
}

type IssueRequest struct {
	ProjectID bson.ObjectId
	CoinType  string
	Address   string
}

type ListRequest struct {
	ProjectID bson.ObjectId
	Limit     uint64
	Page      uint64
	Address   string
}

type service struct {
	db            *Mongo
	addressIssuer *AddressIssuer
}

func NewService(db *Mongo, issuer *AddressIssuer) Service {
	return &service{
		db:            db,
		addressIssuer: issuer,
	}
}

func (s service) Issue(ctx context.Context, issueRequest IssueRequest) (entity.Address, error) {
	// check if project got any configuration about minting
	// if not reject
	var project entity.Project
	err := s.db.ProjectDao.GetByID(issueRequest.ProjectID, &project)

	if err != nil {
		return entity.Address{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project_id": {
				"ADDRESS.ISSUING.NOT_FOUND",
				"Resource unauthenticated or not found",
			},
		}})
	}

	if !common.IsAddress(project.Addresses.MintingAddress) {
		return entity.Address{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"minting_address": {
				"ADDRESS.ISSUING.MALFORM_MINTING_ADDRESS",
				"Malform minting address, double check in project settings",
			},
		}})
	}

	// gen address
	address, accountIndex, err := s.addressIssuer.getAddress(issueRequest.ProjectID, issueRequest.CoinType)
	if err != nil {
		return entity.Address{}, err
	}

	// save issued address, account_index
	newAddress := entity.Address{
		Address:        project.Addresses.MintingAddress,
		ProjectID:      issueRequest.ProjectID,
		AccountIndex:   accountIndex,
		CoinType:       issueRequest.CoinType,
		DepositAddress: address,
	}
	err = s.db.AddressDao.Create(newAddress)

	return newAddress, err
}

func (s service) List(ctx context.Context, listRequest ListRequest) ([]entity.Address, error) {
	return []entity.Address{}, errors.New("Not implemented yet")
}

func (s service) IsIsssueBy(ctx context.Context, address string) (bool, error) {
	return false, errors.New("Not implemented yet")
}
