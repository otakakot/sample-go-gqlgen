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

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/otakakot/sample-go-gqlgen/internal/cache"
	"github.com/otakakot/sample-go-gqlgen/internal/log"
	"github.com/otakakot/sample-go-gqlgen/internal/middleware"
	"github.com/otakakot/sample-go-gqlgen/internal/resolver"
	"github.com/otakakot/sample-go-gqlgen/pkg/graphql"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")

	env := cmp.Or(os.Getenv("ENV"), "production")

	logger := slog.New(log.New(slog.NewJSONHandler(os.Stdout, nil)))
	if env == "local" {
		logger = slog.New(log.New(slog.NewTextHandler(os.Stdout, nil)))
	}

	slog.SetDefault(logger)

	hdl := http.NewServeMux()

	hdl.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	rsl := resolver.New()

	cfg := graphql.Config{
		Resolvers: rsl,
	}

	gql := handler.New(graphql.NewExecutableSchema(cfg))

	gql.AddTransport(transport.POST{})

	gql.Use(extension.AutomaticPersistedQuery{
		Cache: cache.New(),
	})

	hdl.Handle("POST /graphql", middleware.Authorize(gql))

	if env != "production" {
		slog.Info("connect to http://localhost:" + port + "/graphql for GraphQL playground")

		hdl.Handle("GET /graphql", playground.Handler("GraphQL playground", "/graphql"))
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           hdl,
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
