package authdriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	bridgeAuth "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/auth"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) gateway.AuthServiceServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return gateway.AuthServiceKitServer{
		RequestTokenHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.RequestToken,
			decodeRequestLoginTokenGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeRequestLoginTokenGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		LoginHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.Login,
			decodeLoginGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeLoginGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeRequestLoginTokenGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.RequestTokenRequest)

	return RequestTokenRequest{
		Request: bridgeAuth.RqTokenData{
			Address: req.Address,
		},
	}, nil
}

func encodeRequestLoginTokenGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(RequestTokenResponse)

	return &gateway.RequestTokenResponse{
		Token: resp.Token.Token,
	}, nil
}

func decodeLoginGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.AuthServiceLoginRequest)

	return LoginRequest{
		Request: bridgeAuth.Token{
			Address:   req.Address,
			Token:     req.Token,
			Signature: req.Signature,
		},
	}, nil
}

func encodeLoginGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(LoginResponse)

	return &gateway.AuthServiceLoginResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
