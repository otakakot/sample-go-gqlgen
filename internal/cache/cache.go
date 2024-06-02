package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/otakakot/sample-go-gqlgen/internal/domain"
)

var _ graphql.Cache = (*Cache)(nil)

type Cache struct {
	cache *expirable.LRU[string, any]
}

func New() *Cache {
	return &Cache{
		cache: expirable.NewLRU[string, any](256, nil, 600*time.Second),
	}
}

// Add implements graphql.Cache.
func (cc *Cache) Add(ctx context.Context, key string, value interface{}) {
	slog.InfoContext(ctx, "cache add")

	uid := domain.CtxValUserID(ctx)

	if uid == "" {
		return
	}

	// Add user ID to key to avoid collision.
	key = uid + key

	slog.InfoContext(ctx, "cache add", "key", key)

	_ = cc.cache.Add(key, value)
}

// Get implements graphql.Cache.
func (cc *Cache) Get(ctx context.Context, key string) (value interface{}, ok bool) {
	slog.InfoContext(ctx, "cache get")

	uid := domain.CtxValUserID(ctx)

	if uid == "" {
		return nil, false
	}

	// Add user ID to key to avoid collision.
	key = uid + key

	return cc.cache.Get(key)
}
