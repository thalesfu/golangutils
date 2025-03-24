package hertzext

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/thalesfu/golangutils"
	"github.com/thalesfu/golangutils/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

const HertzServiceHandler = "hertz-service-handler"

func ServerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		now := time.Now()
		ctx, logData := logging.InitializeLogContext(ctx)
		serviceHandlerLogData := make(map[string]string)
		logData[HertzServiceHandler] = serviceHandlerLogData

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithTimestamp(now),
			oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		}

		parentSpan := oteltrace.SpanFromContext(ctx)
		tracer := otel.Tracer(parentSpan.SpanContext().TraceID().String())
		ctx, span := tracer.Start(ctx, HertzServiceHandler, opts...)
		defer func() {
			for k, v := range serviceHandlerLogData {
				span.SetAttributes(attribute.String(k, v))
			}
			span.End()

			hlog.CtxInfof(ctx, "Service Handler %s:%s", c.Method(), c.URI().String())
		}()

		serviceHandlerLogData["trace_id"] = span.SpanContext().TraceID().String()
		serviceHandlerLogData["path"] = string(c.Path())
		serviceHandlerLogData["method"] = string(c.Method())
		serviceHandlerLogData["request"] = string(c.Request.Body())

		c.Request.Header.VisitAll(func(k, v []byte) {
			serviceHandlerLogData[fmt.Sprintf("header:%s", k)] = string(v)
		})

		cookies := c.Request.Header.Cookies()

		for _, cookie := range cookies {
			serviceHandlerLogData[fmt.Sprintf("cookie:%s", cookie.Key())] = string(cookie.Value())
		}

		c.QueryArgs().VisitAll(func(k, v []byte) {
			serviceHandlerLogData[fmt.Sprintf("query:%s", k)] = string(v)
		})

		serviceHandlerLogData["time"] = now.String()
		serviceHandlerLogData["client_ip"] = c.ClientIP()
		serviceHandlerLogData["session_id"] = string(c.GetHeader("x-session-id"))

		serviceHandlerLogData["host"] = golangutils.GetHostname()
		serviceHandlerLogData["ip"] = golangutils.GetIP()

		c.Next(ctx)

		serviceHandlerLogData["status_code"] = strconv.Itoa(c.Response.StatusCode())
		serviceHandlerLogData["response"] = string(c.Response.Body())
	}
}

func ErrorMiddleware() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(PanicHandler))
}

func PanicHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	currentSpan := oteltrace.SpanFromContext(ctx)

	logContext := logging.GetLogContext(ctx)
	m := logContext[HertzServiceHandler]
	m["error"] = fmt.Sprint(err)
	m["error_stack"] = string(stack)

	currentSpan.SetStatus(codes.Error, "panic occurred")
	currentSpan.RecordError(errors.New(fmt.Sprint(err)), oteltrace.WithStackTrace(true))

	c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err), "stack": string(stack)})
}
