package Middleware

import (
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	"github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func AdminAuthMiddleware(tokenMaker Token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(Utils.AuthorizationHeaderKey)
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
		if authorizationType != Utils.AuthorizationTypeBearer {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization type"})
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		userRole := payload.Role
		if userRole != "Admin" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization role"})
			return
		}
		c.Set(Utils.AuthorizationPayloadKey, payload)
		c.Next()
	}
}
