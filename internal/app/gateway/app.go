package mga

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	appkitendpoint "github.com/sagikazarmark/appkit/endpoint"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxtransport "github.com/sagikazarmark/kitx/transport"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	"google.golang.org/grpc"

	"github.com/anhntbk08/gateway/internal/app/gateway/httpbin"
	"github.com/anhntbk08/gateway/internal/app/gateway/landing/landingdriver"

	// "github.com/anhntbk08/gateway/internal/app/gateway/bridge/adapter"
	// "github.com/anhntbk08/gateway/internal/app/gateway/bridge/adapter/ent"
	// "github.com/anhntbk08/gateway/internal/app/gateway/bridge/adapter/ent/migrate"

	// TODO find way to merge all small services part into 1 sub-service with driver, store adaptor ...
	authv1 "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	bridgeAuth "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/auth"
	bridgeAuthDriver "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/auth/authdriver"
	// todov1beta1 "github.com/anhntbk08/gateway/.gen/api/proto/todo/v1beta1"
	// "github.com/anhntbk08/gateway/internal/app/gateway/todo"
	// "github.com/anhntbk08/gateway/internal/app/gateway/todo/todoadapter"
	// "github.com/anhntbk08/gateway/internal/app/gateway/todo/todoadapter/ent"
	// "github.com/anhntbk08/gateway/internal/app/gateway/todo/todoadapter/ent/migrate"
	// "github.com/anhntbk08/gateway/internal/app/gateway/todo/tododriver"
	// "github.com/anhntbk08/gateway/internal/app/gateway/todo/todogen"
)

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	storage string,
	db *sql.DB,
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

	transportErrorHandler := kitxtransport.NewErrorHandler(errorHandler)

	// httpServerOptions := []kithttp.ServerOption{
	// 	kithttp.ServerErrorHandler(transportErrorHandler),
	// 	kithttp.ServerErrorEncoder(kitxhttp.NewJSONProblemErrorEncoder(appkithttp.NewDefaultProblemConverter())),
	// 	kithttp.ServerBefore(correlation.HTTPToContext(), kithttp.PopulateRequestContext),
	// }

	grpcServerOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transportErrorHandler),
		kitgrpc.ServerBefore(correlation.GRPCToContext()),
	}

	{
		// eventBus, _ := cqrs.NewEventBus(
		// 	publisher,
		// 	func(eventName string) string { return todoTopic },
		// 	cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		// )

		// var store todo.Store = todo.NewInMemoryStore()
		// if storage == "database" {
		// 	client := ent.NewClient(ent.Driver(entsql.OpenDB("mysql", db)))
		// 	err := client.Schema.Create(
		// 		context.Background(),
		// 		migrate.WithDropIndex(true),
		// 		migrate.WithDropColumn(true),
		// 	)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	store = todoadapter.NewEntStore(client)
		// }

		service := bridgeAuth.NewService(
			ulidgen.NewGenerator(),
		)
		service = bridgeAuthDriver.LoggingMiddleware(logger)(service)
		service = bridgeAuthDriver.InstrumentationMiddleware()(service)

		endpoints := bridgeAuthDriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		// bridgeAuthDriver.RegisterHTTPHandlers(
		// 	endpoints,
		// 	httpRouter.PathPrefix("/bridge").Subrouter(),
		// 	kitxhttp.ServerOptions(httpServerOptions),
		// )

		authv1.RegisterAuthServiceServer(
			grpcServer,
			bridgeAuthDriver.MakeGRPCServer(
				endpoints,
				kitxgrpc.ServerOptions(grpcServerOptions),
			),
		)
		// todov1beta1.RegisterTodoListServer(
		// 	grpcServer,
		// 	bridgeAuthDriver.MakeGRPCServer(
		// 		endpoints,
		// 		kitxgrpc.ServerOptions(grpcServerOptions),
		// 	),
		// )

	}

	landingdriver.RegisterHTTPHandlers(httpRouter)
	httpRouter.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(logger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))
}

// // RegisterEventHandlers registers event handlers in a message router.
// func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger Logger) error {
// 	todoEventProcessor, _ := cqrs.NewEventProcessor(
// 		[]cqrs.EventHandler{
// 			todogen.NewMarkedAsCompleteEventHandler(todo.NewLogEventHandler(logger), "marked_as_complete"),
// 		},
// 		func(eventName string) string { return todoTopic },
// 		func(handlerName string) (message.Subscriber, error) { return subscriber, nil },
// 		cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
// 		watermilllog.New(logger.WithFields(map[string]interface{}{"component": "watermill"})),
// 	)

// 	err := todoEventProcessor.AddHandlersToRouter(router)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
