package projectdriver

import (
	"context"

	bridgev1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/jwt"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
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
		ListHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			VerifyToken(jwtKeyFunc, jwt.SigningMethodHS256, UserClaimFactory)(endpoints.List),
			decodeListGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeListGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		UpdateHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			VerifyToken(jwtKeyFunc, jwt.SigningMethodHS256, UserClaimFactory)(endpoints.Update),
			decodeUpdateGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeUpdateGRPCResponse, errorEncoder),
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
		Id:   resp.Project.ID.String(),
		User: resp.Project.User.Hex(),
		Keys: &bridgev1.Keys{
			Id:     resp.Project.Keys.ID,
			Secret: resp.Project.Keys.Secret,
		},
		Name: resp.Project.Name,
	}, nil
}

func decodeListGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return ListRequest{}, nil
}

func encodeListGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ListResponse)

	projects := make([]*bridgev1.Project, len(resp.Projects))

	for i, t := range resp.Projects {
		projects[i] = &bridgev1.Project{
			Id:   t.ID.String(),
			Name: t.Name,
			User: t.User.Hex(),
			Keys: &bridgev1.Keys{
				Id:     t.Keys.ID,
				Secret: t.Keys.Secret,
			},
			Addresses: &bridgev1.Addresses{},
			Security: &bridgev1.Security{
				WhiteListDomains: t.Security.WhileListAddresses,
				WhiteListIps:     t.Security.WhileListOrigins,
			},
		}
	}

	return &bridgev1.ListResponse{
		Projects: projects,
	}, resp.Err
}

func decodeUpdateGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*bridgev1.UpdateRequest)

	return UpdateRequest{
		Project: entity.Project{
			ID:   bson.ObjectIdHex(req.Id),
			Name: req.Name,
		},
	}, nil
}

func encodeUpdateGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(UpdateResponse)

	success := true
	if resp.Err != nil {
		success = false
	}
	return &bridgev1.UpdateResponse{
		Success: success,
	}, resp.Err
}
