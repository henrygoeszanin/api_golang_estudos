package middlewares

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	jwtv4 "github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/services"
	"github.com/henrygoeszanin/api_golang_estudos/config"
)

type login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SetupJWTMiddleware configura o middleware JWT para autenticação
func SetupJWTMiddleware(userService services.UserService, cfg *config.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "library-api",
		Key:         []byte(cfg.JWTSecret),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24 * 7,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwtv4.MapClaims {
			if v, ok := data.(*dtos.UserResponseDTO); ok {
				return jwtv4.MapClaims{
					"id":       v.ID,
					"email":    v.Email,
					"is_admin": v.IsAdmin,
				}
			}
			return jwtv4.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &dtos.UserResponseDTO{
				ID:      uint(claims["id"].(float64)),
				Email:   claims["email"].(string),
				IsAdmin: claims["is_admin"].(bool),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			user, err := userService.AuthenticateUser(loginVals.Email, loginVals.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &dtos.UserResponseDTO{
				ID:      user.ID,
				Name:    user.Name,
				Email:   user.Email,
				IsAdmin: user.IsAdmin,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true // Implementar lógica de autorização específica
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

// AdminRequired verifica se o usuário é um administrador
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		if isAdmin, exists := claims["is_admin"]; !exists || !isAdmin.(bool) {
			c.JSON(403, gin.H{"error": "acesso restrito a administradores"})
			c.Abort()
			return
		}
		c.Next()
	}
}
