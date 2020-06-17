package authdriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	bridgev1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	bridgeAuth "github.com/anhntbk08/gateway/internal/app/tmbridgev2/bridge/service/auth"
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
		LoginHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.Login,
			decodeLoginGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeLoginGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeRequestLoginTokenGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
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
		Token: resp.Token.Token,
	}, nil
}

func decodeLoginGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*bridgev1.AuthServiceLoginRequest)

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

	return &bridgev1.AuthServiceLoginResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
