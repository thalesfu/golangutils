package hertzext

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/sirupsen/logrus"
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

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithTimestamp(now),
			oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		}

		parentSpan := oteltrace.SpanFromContext(ctx)
		tracer := otel.Tracer(parentSpan.SpanContext().TraceID().String())
		ctx, span := tracer.Start(ctx, HertzServiceHandler, opts...)
		ctx, logStore := logging.InitializeContextLogStore(ctx, HertzServiceHandler)
		defer func() {
			logData := logStore.GetAll()
			for k, v := range logData {
				span.SetAttributes(attribute.String(k, v))
			}
			span.End()

			logrus.WithContext(ctx).Infof("Service Handler %s:%s", c.Method(), c.URI().String())
		}()

		logStore.Set("trace_id", span.SpanContext().TraceID().String())
		logStore.Set("path", string(c.Path()))
		logStore.Set("method", string(c.Method()))
		logStore.Set("request", string(c.Request.Body()))

		c.Request.Header.VisitAll(func(k, v []byte) {
			logStore.Set(fmt.Sprintf("header:%s", k), string(v))
		})

		cookies := c.Request.Header.Cookies()

		for _, cookie := range cookies {
			logStore.Set(fmt.Sprintf("cookie:%s", cookie.Key()), string(cookie.Value()))
		}

		c.QueryArgs().VisitAll(func(k, v []byte) {
			logStore.Set(fmt.Sprintf("query:%s", k), string(v))
		})

		logStore.Set("time", now.String())
		logStore.Set("client_ip", c.ClientIP())
		logStore.Set("session_id", string(c.GetHeader("x-session-id")))

		logStore.Set("host", golangutils.GetHostname())
		logStore.Set("ip", golangutils.GetIP())

		c.Next(ctx)

		logStore.Set("status_code", strconv.Itoa(c.Response.StatusCode()))
		logStore.Set("response", string(c.Response.Body()))
	}
}

func ErrorMiddleware() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(PanicHandler))
}

func PanicHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	currentSpan := oteltrace.SpanFromContext(ctx)

	if logStore, ok := logging.GetContextLogStore(ctx); ok {
		logStore.Set("error", fmt.Sprint(err))
		logStore.Set("error_stack", string(stack))
	}

	currentSpan.SetStatus(codes.Error, "panic occurred")
	currentSpan.RecordError(errors.New(fmt.Sprint(err)), oteltrace.WithStackTrace(true))

	c.JSON(consts.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err), "stack": string(stack)})
}
