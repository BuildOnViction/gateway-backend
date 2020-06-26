package authdriver

import (
	"context"

	"emperror.dev/errors"
	authService "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/auth"
	common "github.com/anhntbk08/gateway/internal/common"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/peer"
)

// Middleware is a service middleware.
type Middleware func(authService.Service) authService.Service

// defaultMiddleware helps implementing partial middleware.
type defaultMiddleware struct {
	service authService.Service
}

var limiter = rate.NewLimiter(10, 3)

func (m defaultMiddleware) RequestToken(ctx context.Context, request authService.RqTokenData) (authService.Token, error) {
	if limiter.Allow() == false {
		return authService.Token{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"request": {
				"TOO_MANY_REQUESTS",
				"Too many request",
			},
		}})
	}
	return m.service.RequestToken(ctx, request)
}

func (m defaultMiddleware) Login(ctx context.Context, request authService.Token) (string, error) {
	return m.service.Login(ctx, request)
}

// LoggingMiddleware is a service level logging middleware.
func LoggingMiddleware(logger authService.Logger) Middleware {
	return func(next authService.Service) authService.Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   authService.Service
	logger authService.Logger
}

func (mw loggingMiddleware) RequestToken(ctx context.Context, request authService.RqTokenData) (authService.Token, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("Request token")

	token, err := mw.next.RequestToken(ctx, request)
	if err != nil {
		return token, err
	}

	logger.Info("Requested token", map[string]interface{}{"token": token.Message})

	return token, err
}

func (mw loggingMiddleware) Login(ctx context.Context, request authService.Token) (string, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info(request.Address + " trying to login in ")

	resp, err := mw.next.Login(ctx, request)
	if err != nil {
		return "", err
	}

	logger.Info("Logged in", map[string]interface{}{"address": request.Address})

	return resp, err
}

// Business metrics
// nolint: gochecknoglobals,lll
var (
	RequestLoginTokenCount = stats.Int64("request_login_token_count", "Number of todo items created", stats.UnitDimensionless)
	LoginCount             = stats.Int64("login_count", "Number of todo items marked complete", stats.UnitDimensionless)
)

// nolint: gochecknoglobals
var (
	RequestLoginTokenCountView = &view.View{
		Name:        "auth.request_login_token_count",
		Description: "Count of number requests for login token",
		Measure:     RequestLoginTokenCount,
		Aggregation: view.Count(),
	}

	LoginCountView = &view.View{
		Name:        "auth.login_count",
		Description: "Count of login request",
		Measure:     LoginCount,
		Aggregation: view.Count(),
	}
)

// InstrumentationMiddleware is a service level instrumentation middleware.
func InstrumentationMiddleware() Middleware {
	return func(next authService.Service) authService.Service {
		return instrumentationMiddleware{
			Service: defaultMiddleware{next},
			next:    next,
		}
	}
}

type instrumentationMiddleware struct {
	authService.Service
	next authService.Service
}

func (mw instrumentationMiddleware) RequestToken(ctx context.Context, request authService.RqTokenData) (authService.Token, error) {
	token, err := mw.next.RequestToken(ctx, request)
	if err != nil {
		return token, err
	}

	p, ok := peer.FromContext(ctx)

	if ok {
		if span := trace.FromContext(ctx); span != nil {
			span.AddAttributes(trace.StringAttribute("ip", p.Addr.String()))
		}
	} else {
		if span := trace.FromContext(ctx); span != nil {
			span.AddAttributes(trace.StringAttribute("address", request.Address))
		}
	}

	stats.Record(ctx, RequestLoginTokenCount.M(1))

	return token, nil
}

func (mw instrumentationMiddleware) Login(ctx context.Context, request authService.Token) (string, error) {
	token, err := mw.next.Login(ctx, request)
	if err != nil {
		return token, err
	}

	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("address", request.Address))
	}

	stats.Record(ctx, LoginCount.M(1))

	return token, nil
}
