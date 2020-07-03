package gateway

import (
	"context"
	"net/http"

	"emperror.dev/emperror"
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
	bus "github.com/anhntbk08/gateway/internal/app/tmbridgev2/bus"
	job "github.com/anhntbk08/gateway/internal/app/tmbridgev2/job"

	bridgeAuth "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/auth"
	bridgeAuthDriver "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/auth/authdriver"

	project "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/project"
	projectDriver "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/project/projectdriver"

	address "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/address"
	addressDriver "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/address/addressdriver"

	store "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/anhntbk08/gateway/internal/platform/database"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
)

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	// publisher message.Publisher,
	dbConfig database.Config,
	jwtConfig common.JWT,
	xPubkeys map[string]string,
	jobQueueConfig common.JobqueueConfig,
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

	// job server run internal
	bus, err := bus.NewBus(jobQueueConfig)
	emperror.Panic(err)
	emperror.Panic(job.StartServer(bus))

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
			bus,
		)
		service = projectDriver.LoggingMiddleware(logger)(service)
		service = projectDriver.InstrumentationMiddleware()(service)

		endpoints := projectDriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		projectServerOptions := append(grpcServerOptions,
			kitgrpc.ServerBefore(gokitjwt.GRPCToContext()),
		)

		gateway.RegisterProjectServiceServer(
			grpcServer,
			projectDriver.MakeGRPCServer(
				endpoints,
				jwtConfig.Key,
				kitxgrpc.ServerOptions(projectServerOptions),
			),
		)
	}

	{
		service := address.NewService(
			mongoConnection,
			address.NewIssuer(mongoConnection, xPubkeys),
		)
		service = addressDriver.LoggingMiddleware(logger)(service)
		// service = addressDriver.InstrumentationMiddleware()(service)

		endpoints := addressDriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		addressServerOptions := append(grpcServerOptions,
			kitgrpc.ServerBefore(gokitjwt.GRPCToContext()),
		)

		gateway.RegisterAddressServiceServer(
			grpcServer,
			addressDriver.MakeGRPCServer(
				endpoints,
				mongoConnection,
				kitxgrpc.ServerOptions(addressServerOptions),
			),
		)
	}

	landingdriver.RegisterHTTPHandlers(httpRouter)
	httpRouter.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(logger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))

}
