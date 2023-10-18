package Middleware

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func SentryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("ip", c.ClientIP())
		}
		c.Next()
	}
}
