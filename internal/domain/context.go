package domain

import "context"

type CtxKey string

const (
	UIDKey CtxKey = "uid"
)

func CtxWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UIDKey, userID)
}

func CtxValUserID(ctx context.Context) string {
	uid, ok := ctx.Value(UIDKey).(string)
	if !ok {
		return ""
	}

	return uid
}
