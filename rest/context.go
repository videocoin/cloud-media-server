package rest

import (
	"context"
)

type key int

const (
	serviceName key = 1
)

func NewContextWithServiceName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, serviceName, name)
}

func ServiceNameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(serviceName).(string)
	return name, ok
}
