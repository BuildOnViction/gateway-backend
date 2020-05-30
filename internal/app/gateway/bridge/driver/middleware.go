package driver

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	bridgeServices "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service"
)

// Middleware is a service middleware.
type Middleware func(bridgeServices.Service) bridgeServices.Service

// defaultMiddleware helps implementing partial middleware.
type defaultMiddleware struct {
	service bridgeServices.Service
}

func (m defaultMiddleware) AddItem(ctx context.Context, newItem bridgeServices.NewItem) (bridgeServices.Item, error) {
	return m.service.AddItem(ctx, newItem)
}

// LoggingMiddleware is a service level logging middleware.
func LoggingMiddleware(logger bridgeServices.Logger) Middleware {
	return func(next bridgeServices.Service) bridgeServices.Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   bridgeServices.Service
	logger bridgeServices.Logger
}

func (mw loggingMiddleware) AddItem(ctx context.Context, newItem bridgeServices.NewItem) (bridgeServices.Item, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("adding item")

	id, err := mw.next.AddItem(ctx, newItem)
	if err != nil {
		return id, err
	}

	logger.Info("added item", map[string]interface{}{"item_id": id})

	return id, err
}

// Business metrics
// nolint: gochecknoglobals,lll
var (
	CreatedTodoItemCount  = stats.Int64("created_todo_item_count", "Number of todo items created", stats.UnitDimensionless)
	CompleteTodoItemCount = stats.Int64("complete_todo_item_count", "Number of todo items marked complete", stats.UnitDimensionless)
)

// nolint: gochecknoglobals
var (
	CreatedTodoItemCountView = &view.View{
		Name:        "todo_item_created_count",
		Description: "Count of todo items created",
		Measure:     CreatedTodoItemCount,
		Aggregation: view.Count(),
	}

	CompleteTodoItemCountView = &view.View{
		Name:        "todo_item_complete_count",
		Description: "Count of todo items complete",
		Measure:     CompleteTodoItemCount,
		Aggregation: view.Count(),
	}
)

// InstrumentationMiddleware is a service level instrumentation middleware.
func InstrumentationMiddleware() Middleware {
	return func(next bridgeServices.Service) bridgeServices.Service {
		return instrumentationMiddleware{
			Service: defaultMiddleware{next},
			next:    next,
		}
	}
}

type instrumentationMiddleware struct {
	bridgeServices.Service
	next bridgeServices.Service
}

func (mw instrumentationMiddleware) AddItem(ctx context.Context, newItem bridgeServices.NewItem) (bridgeServices.Item, error) {
	item, err := mw.next.AddItem(ctx, newItem)
	if err != nil {
		return item, err
	}

	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("item_id", item.ID))
	}

	stats.Record(ctx, CreatedTodoItemCount.M(1))

	return item, nil
}

func (mw instrumentationMiddleware) UpdateItem(ctx context.Context, id string, itemUpdate bridgeServices.ItemUpdate) (bridgeServices.Item, error) { // nolint: lll
	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("item_id", id))
	}

	if itemUpdate.Completed != nil && *itemUpdate.Completed {
		stats.Record(ctx, CompleteTodoItemCount.M(1))
	}

	return mw.next.UpdateItem(ctx, id, itemUpdate)
}
