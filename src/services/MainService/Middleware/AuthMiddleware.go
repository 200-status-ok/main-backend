package Middleware

import (
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware(tokenMaker Token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(utils.AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing authorization header"})
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header"})
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != utils.AuthorizationTypeBearer {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization type"})
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Set(utils.AuthorizationPayloadKey, payload)
		c.Next()
	}
}
