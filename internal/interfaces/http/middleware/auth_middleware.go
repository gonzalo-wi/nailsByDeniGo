package middleware

import (
	"net/http"
	"strings"

	"apiGoShei/internal/infrastructure/security"

	"github.com/gin-gonic/gin"
)

// RequireRole verifica que el rol del token esté dentro de los roles permitidos.
// Debe usarse después de AuthMiddleware.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		raw, exists := c.Get(ClaimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
			return
		}
		claims := raw.(*security.Claims)
		if !allowed[claims.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
			return
		}
		c.Next()
	}
}

const ClaimsKey = "claims"

// AuthMiddleware valida el JWT en el header Authorization: Bearer <token>.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token invalido: Bearer <token>"})
			return
		}
		claims, err := security.ParseToken(parts[1], jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
			return
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}
