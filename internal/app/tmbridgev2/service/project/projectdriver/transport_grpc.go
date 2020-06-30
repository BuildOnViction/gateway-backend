package projectdriver

import (
	"context"

	"emperror.dev/errors"
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/jwt"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, jwtkey string, options ...kitgrpc.ServerOption) gateway.ProjectServiceKitServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())
	jwtKeyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	}

	return gateway.ProjectServiceKitServer{
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
		DeleteHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			VerifyToken(jwtKeyFunc, jwt.SigningMethodHS256, UserClaimFactory)(endpoints.Delete),
			decodeDeleteGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeDeleteGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		GetOneHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			VerifyToken(jwtKeyFunc, jwt.SigningMethodHS256, UserClaimFactory)(endpoints.GetOne),
			decodeGetOneGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeGetOneGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeCreateGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.CreateRequest)

	return CreateRequest{
		Name: req.Name,
	}, nil
}

func encodeCreateGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CreateResponse)

	return &gateway.CreateResponse{
		Id:   resp.Project.ID.Hex(),
		User: resp.Project.User.Hex(),
		Keys: &gateway.Keys{
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

	projects := make([]*gateway.Project, len(resp.Projects))

	for i, t := range resp.Projects {
		projects[i] = &gateway.Project{
			Id:   t.ID.Hex(),
			Name: t.Name,
			User: t.User.Hex(),
			Keys: &gateway.Keys{
				Id:     t.Keys.ID,
				Secret: t.Keys.Secret,
			},
			Addresses: &gateway.Addresses{
				MintingAddress:      t.Addresses.MintingAddress,
				WatchSmartContracts: t.Addresses.WatchSmartContracts,
			},
			Security: &gateway.Security{
				WhiteListDomains: t.Security.WhileListDomains,
				WhiteListIps:     t.Security.WhileListIps,
			},
			Notification: &gateway.Notification{
				WebHook: t.Notification.WebHook,
				Emails:  t.Notification.Emails,
			},
			CreatedAt: t.CreatedAt.Unix(),
			UpdatedAt: t.UpdatedAt.Unix(),
		}
	}

	return &gateway.ListResponse{
		Projects: projects,
	}, resp.Err
}

func decodeUpdateGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.UpdateRequest)

	addresses := entity.ProjectAddresses{}
	if req.Addresses != nil {
		addresses = entity.ProjectAddresses{
			MintingAddress:      req.Addresses.MintingAddress,
			WatchSmartContracts: req.Addresses.WatchSmartContracts,
		}
	}

	notification := entity.Notification{}
	if req.Notification != nil {
		notification = entity.Notification{
			WebHook: req.Notification.WebHook,
			Emails:  req.Notification.Emails,
		}
	}

	security := entity.Security{}
	if req.Security != nil {
		security = entity.Security{
			WhileListDomains: req.Security.WhiteListDomains,
			WhileListIps:     req.Security.WhiteListIps,
		}
	}

	if common.IsValidMongoID(req.Id) {
		return UpdateRequest{
			Project: entity.Project{
				ID:           bson.ObjectIdHex(req.Id),
				Name:         req.Name,
				Addresses:    addresses,
				Notification: notification,
				Security:     security,
				Status:       req.Status,
			},
		}, nil
	} else {
		return UpdateRequest{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.UPDATING.MALFORMED_PROJECT_ID",
				"Malformed project id",
			},
		}})
	}
}

func encodeUpdateGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(UpdateResponse)

	success := true
	if resp.Err != nil {
		success = false
	}
	return &gateway.UpdateResponse{
		Success: success,
	}, resp.Err
}

func decodeDeleteGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.DeleteRequest)
	if common.IsValidMongoID(req.Id) {
		return DeleteRequest{
			Id: bson.ObjectIdHex(req.Id),
		}, nil
	} else {
		return DeleteRequest{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.DELETING.MALFORMED_PROJECT_ID",
				"Malformed project id",
			},
		}})
	}
}

func encodeDeleteGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(DeleteResponse)

	success := true
	if resp.Err != nil {
		success = false
	}
	return &gateway.DeleteResponse{
		Success: success,
	}, resp.Err
}

func decodeGetOneGRPCRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*gateway.GetOneRequest)
	if common.IsValidMongoID(req.Id) {
		return GetOneRequest{
			Id: bson.ObjectIdHex(req.Id),
		}, nil
	} else {
		return GetOneRequest{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.GETONE.MALFORMED_PROJECT_ID",
				"Malformed project id",
			},
		}})
	}
}

func encodeGetOneGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetOneResponse)

	return &gateway.GetOneResponse{
		Project: &gateway.Project{
			Id:   resp.Project.ID.Hex(),
			Name: resp.Project.Name,
			User: resp.Project.User.Hex(),
			Keys: &gateway.Keys{
				Id:     resp.Project.Keys.ID,
				Secret: resp.Project.Keys.Secret,
			},
			Addresses: &gateway.Addresses{
				MintingAddress:      resp.Project.Addresses.MintingAddress,
				WatchSmartContracts: resp.Project.Addresses.WatchSmartContracts,
			},
			Security: &gateway.Security{
				WhiteListDomains: resp.Project.Security.WhileListDomains,
				WhiteListIps:     resp.Project.Security.WhileListIps,
			},
			Notification: &gateway.Notification{
				WebHook: resp.Project.Notification.WebHook,
				Emails:  resp.Project.Notification.Emails,
			},
			CreatedAt: resp.Project.CreatedAt.Unix(),
			UpdatedAt: resp.Project.UpdatedAt.Unix(),
		},
	}, resp.Err
}
