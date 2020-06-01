package authdriver

import (
	"context"
	"fmt"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	bridgev1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	bridgeAuth "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/auth"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) bridgev1.AuthServiceServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return bridgev1.AuthServiceKitServer{
		RequestTokenHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.RequestToken,
			decodeRequestLoginTokenGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeRequestLoginTokenGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeRequestLoginTokenGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	fmt.Println("bridgev1.AuthServiceLoginRequest ", request)
	req := request.(*bridgev1.RequestTokenRequest)

	return RequestTokenRequest{
		Request: bridgeAuth.RqTokenData{
			Address: req.Address,
		},
	}, nil
}

func encodeRequestLoginTokenGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(RequestTokenResponse)

	return &bridgev1.RequestTokenResponse{
		Token: resp.Token.IssuedToken,
	}, nil
}
