package contexts

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
		return "", errors.New("error to get store id from context")
	}
	return tokenStr, nil
}
