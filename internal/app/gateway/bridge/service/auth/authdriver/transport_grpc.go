package authdriver

import (
	"context"

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
	req := request.(*bridgev1.AuthServiceLoginRequest)

	return bridgev1.AuthServiceLoginRequest{
		Address: req.Address,
	}, nil
}

func encodeRequestLoginTokenGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(bridgeAuth.Token)

	return &bridgev1.RequestTokenResponse{
		Token: resp.IssuedToken,
	}, nil
}

// func marshalItemGRPC(item bridge.Token) *bridgev1.TodoItem {
// 	return &bridgev1.TodoItem{
// 		Id:        item.ID,
// 		Title:     item.Title,
// 		Completed: item.Completed,
// 		Order:     int32(item.Order),
// 	}
// }
