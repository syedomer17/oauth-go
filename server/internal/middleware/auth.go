package middleware

import (
	"fmt"
	"net/http"
	"oauth/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired validates the access_token HttpOnly cookie on every protected route.
// It rejects missing, tampered, or expired tokens before the handler runs.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the short-lived access token set by the OAuth callback.
		rawToken, err := c.Cookie("access_token")
		if err != nil || rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		// Parse and verify the JWT signature using the shared secret.
		// The key func also guards against algorithm-confusion attacks (alg:none / RSA→HMAC).
		token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// Propagate the verified userID so downstream handlers skip re-parsing the token.
		c.Set("userID", userID)
		c.Next()
	}
}
