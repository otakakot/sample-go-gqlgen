package log

import (
	"context"
	"log/slog"

	"github.com/otakakot/sample-go-gqlgen/internal/domain"
)

var keys = []domain.CtxKey{domain.UIDKey}

type Handler struct {
	slog.Handler
}

func New(
	handler slog.Handler,
) *Handler {
	return &Handler{
		Handler: handler,
	}
}

func (hdl *Handler) Handle(
	ctx context.Context,
	r slog.Record,
) error {
	for _, key := range keys {
		if v := ctx.Value(key); v != nil {
			r.AddAttrs(slog.Attr{
				Key:   string(key),
				Value: slog.AnyValue(v),
			})
		}
	}

	return hdl.Handler.Handle(ctx, r)
}
