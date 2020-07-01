package addressdriver

import (
	"context"

	"emperror.dev/errors"
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"

	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/jwt"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/address"
	store "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/globalsign/mgo/bson"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, mongo *store.Mongo, options ...kitgrpc.ServerOption) gateway.AddressServiceServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return gateway.AddressServiceKitServer{
		IssueHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			VerifyAPIKey(mongo)(endpoints.Issue),
			decodeIssueGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeIssueGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeIssueGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.IssueRequest)

	if common.IsValidMongoID(req.ApiToken) {
		return IssueRequest{
			IssueRequest: address.IssueRequest{
				Address:   req.Address,
				CoinType:  req.CoinType,
				ProjectID: bson.ObjectIdHex(req.ApiToken),
			},
		}, nil
	} else {
		return IssueRequest{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project_id": {
				"ADDRESS.ISSUING.MALFORMED_PROJECT_ID",
				"Malformed project id",
			},
		}})
	}
}

func encodeIssueGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(IssueResponse)

	return &gateway.IssueResponse{
		CoinType:       resp.R0.CoinType,
		AccountIndex:   resp.R0.AccountIndex,
		DepositAddress: resp.R0.DepositAddress,
		Address:        resp.R0.Address,
	}, nil
}
