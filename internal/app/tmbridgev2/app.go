package gateway

import (
	"context"
	"net/http"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/gorilla/mux"
	appkitendpoint "github.com/sagikazarmark/appkit/endpoint"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxtransport "github.com/sagikazarmark/kitx/transport"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
	"google.golang.org/grpc"

	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/httpbin"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/landing/landingdriver"

	// TODO find way to merge all small services part into 1 sub-service with driver, store adaptor ...
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	bridgeAuth "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/auth"
	bridgeAuthDriver "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/auth/authdriver"

	project "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/project"
	projectDriver "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/project/projectdriver"

	store "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/anhntbk08/gateway/internal/platform/database"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
)

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	dbConfig database.Config,
	jwtConfig common.JWT,
	logger Logger,
	errorHandler ErrorHandler,
) {
	endpointMiddleware := []endpoint.Middleware{
		correlation.Middleware(),
		opencensus.TraceEndpoint("", opencensus.WithSpanName(func(ctx context.Context, _ string) string {
			name, _ := kitxendpoint.OperationName(ctx)

			return name
		})),
		appkitendpoint.LoggingMiddleware(logger),
	}

	mongoConnection, err := store.NewMongo(dbConfig.Uri, dbConfig.DbName)
	emperror.Panic(err)

	transportErrorHandler := kitxtransport.NewErrorHandler(errorHandler)

	grpcServerOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transportErrorHandler),
	}

	{
		service := bridgeAuth.NewService(
			mongoConnection,
			jwtConfig.Key,
		)
		service = bridgeAuthDriver.LoggingMiddleware(logger)(service)
		service = bridgeAuthDriver.InstrumentationMiddleware()(service)

		endpoints := bridgeAuthDriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		gateway.RegisterAuthServiceServer(
			grpcServer,
			bridgeAuthDriver.MakeGRPCServer(
				endpoints,
				kitxgrpc.ServerOptions(grpcServerOptions),
			),
		)
	}

	{
		service := project.NewService(
			mongoConnection,
		)
		service = projectDriver.LoggingMiddleware(logger)(service)
		service = projectDriver.InstrumentationMiddleware()(service)

		endpoints := projectDriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		grpcServerOptions = append(grpcServerOptions,
			kitgrpc.ServerBefore(gokitjwt.GRPCToContext()),
		)

		gateway.RegisterProjectServiceServer(
			grpcServer,
			projectDriver.MakeGRPCServer(
				endpoints,
				jwtConfig.Key,
				kitxgrpc.ServerOptions(grpcServerOptions),
			),
		)
	}

	landingdriver.RegisterHTTPHandlers(httpRouter)
	httpRouter.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(logger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))
}
