package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/ctxutil"
)

func RequestID() echo.MiddlewareFunc {
	return echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Skipper: func(c *echo.Context) bool {
			return c.Request().Method == http.MethodOptions
		},
		Generator: func() string {
			return uuid.NewString()
		},
		RequestIDHandler: func(c *echo.Context, requestID string) {
			c.Set("request_id", requestID)
			ctx := ctxutil.SetRequestID(c.Request().Context(), requestID)
			c.SetRequest(c.Request().WithContext(ctx))
		},
		TargetHeader: echo.HeaderXRequestID,
	})
}
