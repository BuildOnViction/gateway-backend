package projectdriver

import (
	"context"

	bridgev1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	. "github.com/anhntbk08/gateway/internal/app/gateway/jwt"
	"github.com/dgrijalva/jwt-go"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, jwtkey string, options ...kitgrpc.ServerOption) bridgev1.ProjectServiceKitServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())
	jwtKeyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	}

	return bridgev1.ProjectServiceKitServer{
		CreateHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			VerifyToken(jwtKeyFunc, jwt.SigningMethodHS256, UserClaimFactory)(endpoints.Create),
			decodeCreateGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeCreateGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeCreateGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*bridgev1.CreateRequest)

	return CreateRequest{
		Name: req.Name,
	}, nil
}

func encodeCreateGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CreateResponse)

	return &bridgev1.CreateResponse{
		Id: resp.Project.ID.String(),
		Keys: &bridgev1.Keys{
			Id:     resp.Project.Keys.ID,
			Secret: resp.Project.Keys.Secret,
		},
		Name: resp.Project.Name,
	}, nil
}
