package logging

import "context"

const LogContextKey = "log"

func InitializeLogContext(ctx context.Context) (context.Context, map[string]map[string]string) {
	logger := make(map[string]map[string]string)
	ctx = context.WithValue(ctx, LogContextKey, logger)
	return ctx, logger
}

func GetLogContext(ctx context.Context) map[string]map[string]string {
	if ctx == nil {
		return nil
	}

	return ctx.Value(LogContextKey).(map[string]map[string]string)
}
