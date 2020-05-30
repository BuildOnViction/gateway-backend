package driver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	bridgev1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	bridge "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) bridgev1.AuthServiceServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return bridgev1.AuthServiceKitServer{
		AddItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.AddItem,
			decodeAddItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeAddItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeAddItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*bridgev1.AuthServiceLoginRequest)

	return bridgev1.AuthServiceLoginRequest{
		Address: "0x132131346546464",
	}, nil
}

func encodeAddItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(AddItemResponse)

	return &bridgev1.AddItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func marshalItemGRPC(item bridge.Item) *bridgev1.TodoItem {
	return &bridgev1.TodoItem{
		Id:        item.ID,
		Title:     item.Title,
		Completed: item.Completed,
		Order:     int32(item.Order),
	}
}
