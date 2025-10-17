package clients

import (
	"context"
	"strings"
)

type contextKey string

const authTokenContextKey contextKey = "shopServiceAuthToken"

// ContextWithAuthToken stores the Authorization header for outbound shop-service calls.
func ContextWithAuthToken(ctx context.Context, header string) context.Context {
	header = strings.TrimSpace(header)
	if header == "" {
		return ctx
	}
	return context.WithValue(ctx, authTokenContextKey, header)
}

// AuthTokenFromContext retrieves the Authorization header stored by ContextWithAuthToken.
func AuthTokenFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if token, ok := ctx.Value(authTokenContextKey).(string); ok {
		return token
	}
	return ""
}
