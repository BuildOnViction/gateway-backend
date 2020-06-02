package mga

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

	"github.com/anhntbk08/gateway/internal/app/gateway/httpbin"
	"github.com/anhntbk08/gateway/internal/app/gateway/landing/landingdriver"

	// TODO find way to merge all small services part into 1 sub-service with driver, store adaptor ...
	authv1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	bridgeAuth "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/auth"
	bridgeAuthDriver "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/auth/authdriver"
	store "github.com/anhntbk08/gateway/internal/app/gateway/store"
)

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	storage string,
	dbURI string,
	dbName string,
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

	mongoConnection, err := store.NewMongo(dbURI, dbName)
	emperror.Panic(err)

	transportErrorHandler := kitxtransport.NewErrorHandler(errorHandler)

	grpcServerOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transportErrorHandler),
		kitgrpc.ServerBefore(correlation.GRPCToContext()),
	}

	{
		service := bridgeAuth.NewService(
			mongoConnection,
		)
		service = bridgeAuthDriver.LoggingMiddleware(logger)(service)
		service = bridgeAuthDriver.InstrumentationMiddleware()(service)

		endpoints := bridgeAuthDriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		authv1.RegisterAuthServiceServer(
			grpcServer,
			bridgeAuthDriver.MakeGRPCServer(
				endpoints,
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
