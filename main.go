package main

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/otakakot/sample-go-gqlgen/internal/cache"
	"github.com/otakakot/sample-go-gqlgen/internal/log"
	"github.com/otakakot/sample-go-gqlgen/internal/middleware"
	"github.com/otakakot/sample-go-gqlgen/internal/resolver"
	"github.com/otakakot/sample-go-gqlgen/pkg/gql"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")

	env := cmp.Or(os.Getenv("ENV"), "production")

	logger := slog.New(log.New(slog.NewJSONHandler(os.Stdout, nil)))
	if env == "local" {
		logger = slog.New(log.New(slog.NewTextHandler(os.Stdout, nil)))
	}

	slog.SetDefault(logger)

	mux := http.NewServeMux()

	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	rsl := resolver.New()

	cfg := gql.Config{
		Resolvers: rsl,
	}

	hdl := handler.New(gql.NewExecutableSchema(cfg))

	hdl.AddTransport(transport.POST{})

	hdl.Use(extension.AutomaticPersistedQuery{
		Cache: cache.New(),
	})

	hdl.Use(extension.FixedComplexityLimit(5))

	hdl.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		return next(ctx)
	})

	hdl.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		return next(ctx)
	})

	hdl.AroundRootFields(func(ctx context.Context, next graphql.RootResolver) graphql.Marshaler {
		return next(ctx)
	})

	hdl.AroundFields(func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		return next(ctx)
	})

	mux.Handle("POST /graphql", middleware.Authorize(hdl))

	if env != "production" {
		slog.Info("connect to http://localhost:" + port + "/graphql for GraphQL playground")

		hdl.Use(extension.Introspection{})

		mux.Handle("GET /graphql", playground.Handler("GraphQL playground", "/graphql"))
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		slog.Info("start server listen")

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("start server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	slog.Info("done server shutdown")
}
