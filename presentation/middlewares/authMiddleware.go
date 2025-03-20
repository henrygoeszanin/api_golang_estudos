package middlewares

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	jwtv4 "github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/services"
	"github.com/henrygoeszanin/api_golang_estudos/config"
)

// Estrutura que define os campos necessários para o login
type login struct {
	Email    string `json:"email" binding:"required,email"` // Email do usuário (obrigatório e em formato válido)
	Password string `json:"password" binding:"required"`    // Senha do usuário (obrigatório)
}

// SetupJWTMiddleware configura o middleware JWT para autenticação
//
// Este middleware gerencia todo o processo de autenticação, incluindo:
// - Geração de tokens JWT após login bem-sucedido
// - Validação de tokens em requisições subsequentes
// - Configuração de cookies para armazenamento do token
// - Extração de informações do usuário a partir do token
func SetupJWTMiddleware(userService services.UserService, cfg *config.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		// Configurações básicas do JWT
		Realm:       "library-api",         // Nome do domínio de segurança
		Key:         []byte(cfg.JWTSecret), // Chave secreta para assinar tokens
		Timeout:     time.Hour * 24,        // Tempo de validade do token (24 horas)
		MaxRefresh:  time.Hour * 24 * 7,    // Tempo máximo para renovação do token (7 dias)
		IdentityKey: "id",                  // Campo usado como identificador principal

		// Configurações do cookie que armazenará o token JWT
		SendCookie:     true,                 // Ativa o envio do token como cookie
		CookieName:     "jwt-token",          // Nome do cookie que armazenará o token
		CookieMaxAge:   24 * time.Hour,       // Tempo de vida do cookie (24 horas)
		CookieDomain:   "",                   // Domínio do cookie (vazio = domínio atual)
		SecureCookie:   false,                // Em desenvolvimento = false, em produção com HTTPS = true
		CookieHTTPOnly: true,                 // Impede acesso ao cookie via JavaScript (segurança)
		CookieSameSite: http.SameSiteLaxMode, // Comportamento do cookie em requisições cross-site

		// Configurações de busca do token nas requisições
		// A ordem define a prioridade de busca (primeiro cookie, depois header, por fim query)
		TokenLookup:   "cookie:jwt-token;header:Authorization;query:token",
		TokenHeadName: "Bearer", // Prefixo usado no header Authorization
		TimeFunc:      time.Now, // Função para obter o tempo atual (usada na validação)

		// PayloadFunc define quais informações do usuário serão armazenadas no token JWT
		// Estas informações estarão disponíveis em todas as requisições autenticadas
		PayloadFunc: func(data interface{}) jwtv4.MapClaims {
			if v, ok := data.(*dtos.UserResponseDTO); ok {
				return jwtv4.MapClaims{
					"id":       v.ID,      // ID do usuário
					"email":    v.Email,   // Email do usuário
					"is_admin": v.IsAdmin, // Flag que indica se o usuário é administrador
				}
			}
			return jwtv4.MapClaims{}
		},

		// IdentityHandler extrai as informações do usuário a partir do token JWT
		// Isso permite que os handlers da API acessem os dados do usuário autenticado
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &dtos.UserResponseDTO{
				ID:      uint(claims["id"].(float64)), // Converte o ID de float64 para uint
				Email:   claims["email"].(string),     // Obtém o email
				IsAdmin: claims["is_admin"].(bool),    // Obtém o status de administrador
			}
		},

		// Authenticator é chamado durante o login para validar as credenciais do usuário
		// Esta função recebe email/senha e retorna os dados do usuário se autenticado
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// Extrai email e senha do corpo da requisição
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			// Verifica as credenciais com o serviço de usuário
			user, err := userService.AuthenticateUser(loginVals.Email, loginVals.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			// Retorna os dados do usuário que serão incluídos no token
			return &dtos.UserResponseDTO{
				ID:      user.ID,
				Name:    user.Name,
				Email:   user.Email,
				IsAdmin: user.IsAdmin,
			}, nil
		},

		// Authorizator é chamado em cada requisição protegida para verificar permissões
		// Aqui apenas retornamos true, mas poderia implementar regras mais complexas
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true // Permite acesso a qualquer usuário autenticado
		},

		// Unauthorized define como responder quando a autenticação falha
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
}

// AdminRequired verifica se o usuário autenticado possui permissões de administrador
//
// Este middleware deve ser usado após o middleware JWT principal, para rotas
// que exigem privilégios administrativos, como gerenciamento de usuários e livros.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims (informações) do token JWT
		claims := jwt.ExtractClaims(c)

		// Verifica se a claim "is_admin" existe e se é verdadeira
		if isAdmin, exists := claims["is_admin"]; !exists || !isAdmin.(bool) {
			// Se não for admin, retorna erro 403 (Forbidden)
			c.JSON(403, gin.H{"error": "acesso restrito a administradores"})
			c.Abort() // Interrompe o processamento da requisição
			return
		}

		// Se for admin, continua para o próximo handler
		c.Next()
	}
}
