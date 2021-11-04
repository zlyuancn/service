package grpc

import (
	"context"

	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
)

type Session struct {
	core.ILogger
}

func newSessionFromContext(ctx context.Context) *Session {
	return &Session{
		ILogger: logger.Log,
	}
}

type contextKey struct{}

var sessionContextKey = &contextKey{}

// 从标准context中获取session
func GetSession(ctx context.Context) *Session {
	value := ctx.Value(sessionContextKey)
	se, ok := value.(*Session)
	if !ok {
		se = newSessionFromContext(ctx)
	}
	return se
}
