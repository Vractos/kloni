package contexttools

import (
	"context"
	"errors"
)

type contextKey string

func (c contextKey) String() string {
	return "myContexts" + string(c)
}

var (
	ContextKeyStoreId = contextKey("store_id")
)

func RetrieveStoreIDFromCtx(ctx context.Context) (string, error) {
	tokenStr := ctx.Value(ContextKeyStoreId).(string)
	if tokenStr == "" {
		return "", errors.New("context doesn't have the store id")
	}
	return tokenStr, nil
}
