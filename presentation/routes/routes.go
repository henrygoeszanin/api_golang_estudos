package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/henrygoeszanin/api_golang_estudos/application/services"
	"github.com/henrygoeszanin/api_golang_estudos/config"
	"github.com/henrygoeszanin/api_golang_estudos/infrastructure/repositories"
	"github.com/henrygoeszanin/api_golang_estudos/presentation/handlers"
	"github.com/henrygoeszanin/api_golang_estudos/presentation/middlewares"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Inicializar repositórios
	userRepository := repositories.NewUserRepository(db)

	// Inicializar serviços
	userService := services.NewUserService(userRepository)

	// Configurar middleware JWT
	authMiddleware, err := middlewares.SetupJWTMiddleware(userService, cfg)
	if err != nil {
		panic("JWT middleware setup failed: " + err.Error())
	}

	// Inicializar handlers
	userHandler := handlers.NewUserHandler(userService)

	// Definir grupo base da API
	api := router.Group("/api")

	// Configurar grupos de rotas por domínio
	setupHealthRoutes(api)
	setupAuthRoutes(api, userHandler, authMiddleware)
	setupBookRoutes(api, authMiddleware)
	setupLoanRoutes(api, authMiddleware)
	setupUserRoutes(api, userHandler, authMiddleware)
}

// setupHealthRoutes configura rotas de health check
func setupHealthRoutes(router *gin.RouterGroup) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API funcionando corretamente",
		})
	})
}

// setupAuthRoutes configura rotas de autenticação
func setupAuthRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authMiddleware *jwt.GinJWTMiddleware) {
	auth := router.Group("/auth")
	{
		// Rotas públicas de autenticação
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.GET("/refresh", authMiddleware.RefreshHandler)
		auth.POST("/register", userHandler.Register)
	}
}

// setupBookRoutes configura rotas relacionadas a livros
func setupBookRoutes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Rotas públicas (consulta)
	books := router.Group("/books")
	{
		books.GET("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Listagem de livros não implementada"})
		})

		books.GET("/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Detalhes do livro não implementados"})
		})
	}

	// Rotas administrativas (gerenciamento)
	adminBooks := router.Group("/admin/books")
	adminBooks.Use(authMiddleware.MiddlewareFunc(), middlewares.AdminRequired())
	{
		adminBooks.POST("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Adição de livro não implementada"})
		})

		adminBooks.PUT("/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Atualização de livro não implementada"})
		})

		adminBooks.DELETE("/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Remoção de livro não implementada"})
		})
	}
}

// setupLoanRoutes configura rotas relacionadas a empréstimos
func setupLoanRoutes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Todas as rotas de empréstimos requerem autenticação
	loans := router.Group("/loans")
	loans.Use(authMiddleware.MiddlewareFunc())
	{
		loans.GET("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Listagem de empréstimos não implementada"})
		})

		loans.POST("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Criação de empréstimo não implementada"})
		})

		loans.GET("/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Detalhes do empréstimo não implementados"})
		})

		loans.PUT("/:id/return", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Devolução de empréstimo não implementada"})
		})
	}
}

// setupUserRoutes configura rotas relacionadas a usuários
func setupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Rotas de usuário que precisam de autenticação
	users := router.Group("/users")
	users.Use(authMiddleware.MiddlewareFunc())
	{
		users.GET("/me", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Perfil do usuário não implementado"})
		})

		users.PUT("/me", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Atualização do perfil não implementada"})
		})
	}

	// Rotas administrativas para gerenciamento de usuários
	adminUsers := router.Group("/admin/users")
	adminUsers.Use(authMiddleware.MiddlewareFunc(), middlewares.AdminRequired())
	{
		adminUsers.GET("/", userHandler.List)
		adminUsers.GET("/:id", userHandler.GetByID)
		adminUsers.PUT("/:id", userHandler.Update)
		adminUsers.DELETE("/:id", userHandler.Delete)
		adminUsers.PUT("/:id/promote", userHandler.PromoteToAdmin)
	}
}
