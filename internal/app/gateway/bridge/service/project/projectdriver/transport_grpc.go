package projectdriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	bridgev1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) bridgev1.ProjectServiceKitServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return bridgev1.ProjectServiceKitServer{
		CreateHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.Create,
			decodeCreateGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeCreateGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeCreateGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*bridgev1.CreateRequest)

	return CreateRequest{
		Name: req.Name,
		User: "testing",
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
