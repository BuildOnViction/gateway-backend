package addressdriver

import (
	"context"

	addressService "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/address"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
)

// Middleware is a service middleware.
type Middleware func(addressService.Service) addressService.Service

// defaultMiddleware helps implementing partial middleware.
type defaultMiddleware struct {
	service addressService.Service
}

func (m defaultMiddleware) Issue(ctx context.Context, request addressService.IssueRequest) (entity.Address, error) {
	return m.service.Issue(ctx, request)
}

func (m defaultMiddleware) IsIssueBy(ctx context.Context, address string) (bool, error) {
	return m.service.IsIsssueBy(ctx, address)
}

func (m defaultMiddleware) List(ctx context.Context, request addressService.ListRequest) ([]entity.Address, error) {
	return m.service.List(ctx, request)
}

// LoggingMiddleware is a service level logging middleware.
func LoggingMiddleware(logger addressService.Logger) Middleware {
	return func(next addressService.Service) addressService.Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   addressService.Service
	logger addressService.Logger
}

func (mw loggingMiddleware) Issue(ctx context.Context, request addressService.IssueRequest) (entity.Address, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("Issuing adddress")

	address, err := mw.next.Issue(ctx, request)
	if err != nil {
		return address, err
	}

	logger.Info("Issued address", map[string]interface{}{
		"projectId":       address.ProjectID,
		"deposit_address": address.DepositAddress,
		"cointype":        address.CoinType})

	return address, err
}

func (mw loggingMiddleware) List(ctx context.Context, request addressService.ListRequest) ([]entity.Address, error) {
	return mw.next.List(ctx, request)
}

func (mw loggingMiddleware) IsIsssueBy(ctx context.Context, address string) (bool, error) {
	return mw.next.IsIsssueBy(ctx, address)
}
