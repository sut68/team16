package middleware

import (
	"backend/config"
	"backend/entity"
	"backend/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ""
		if authHeader := ctx.GetHeader("Authorization"); authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}
		if tokenStr == "" {
			if cookie, err := ctx.Request.Cookie("access_token"); err == nil {
				tokenStr = cookie.Value
			}
		}
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
			return 
		}
		claims, err := services.ValidateJWT(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return 
		}

		var user entity.User
		if err := config.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return 
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

// double-submit cookie check for unsafe methods
func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			ctx.Next()
			return 
		}

		cookie, err := ctx.Request.Cookie("csrf_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing csrf cookie"})
			return 
		}
		csrfHeader := ctx.GetHeader("x-CSRF-Token")
		if csrfHeader == "" || csrfHeader != cookie.Value {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "csrf token mismatch"})
			return 
		}
		ctx.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
  return func(c *gin.Context) {
    u, ok := c.Get("currentUser")
    if !ok {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"unauthenticated"})
      return
    }
    user := u.(entity.User)
    if user.Role == nil || user.Role.Name != "admin" {
      c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"admin required"})
      return
    }
    c.Next()
  }
}