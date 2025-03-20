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
	bookRepository := repositories.NewBookRepository(db)
	loanRepository := repositories.NewLoanRepository(db)

	// Inicializar serviços
	userService := services.NewUserService(userRepository)
	bookService := services.NewBookService(bookRepository)
	loanService := services.NewLoanService(loanRepository, bookRepository)

	// Configurar middleware JWT
	authMiddleware, err := middlewares.SetupJWTMiddleware(userService, cfg)
	if err != nil {
		panic("JWT middleware setup failed: " + err.Error())
	}

	// Inicializar handlers
	userHandler := handlers.NewUserHandler(userService)
	bookHandler := handlers.NewBookHandler(bookService)
	loanHandler := handlers.NewLoanHandler(loanService)

	// Definir grupo base da API
	api := router.Group("/api")

	// Configurar grupos de rotas por domínio
	setupHealthRoutes(api)
	setupAuthRoutes(api, userHandler, authMiddleware)
	setupBookRoutes(api, bookHandler, authMiddleware)
	setupLoanRoutes(api, loanHandler, authMiddleware)
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
func setupBookRoutes(router *gin.RouterGroup, bookHandler *handlers.BookHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Rotas públicas (consulta)
	books := router.Group("/books")
	{
		books.GET("/", bookHandler.List)
		books.GET("/:id", bookHandler.GetByID)
	}

	// Rotas administrativas (gerenciamento)
	adminBooks := router.Group("/admin/books")
	adminBooks.Use(middlewares.TokenExtractor(), authMiddleware.MiddlewareFunc(), middlewares.AdminRequired())
	{
		adminBooks.POST("/", bookHandler.Create)
		adminBooks.PUT("/:id", bookHandler.Update)
		adminBooks.DELETE("/:id", bookHandler.Delete)
	}
}

// setupLoanRoutes configura rotas relacionadas a empréstimos
func setupLoanRoutes(router *gin.RouterGroup, loanHandler *handlers.LoanHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Todas as rotas de empréstimos requerem autenticação
	loans := router.Group("/loans")
	loans.Use(middlewares.TokenExtractor(), authMiddleware.MiddlewareFunc())
	{
		loans.GET("/", loanHandler.List)
		loans.POST("/", loanHandler.Create)
		loans.GET("/:id", loanHandler.GetByID)
		loans.PUT("/:id/return", loanHandler.ReturnLoan)
	}
}

// setupUserRoutes configura rotas relacionadas a usuários
func setupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authMiddleware *jwt.GinJWTMiddleware) {
	// Rotas de usuário que precisam de autenticação
	users := router.Group("/users")
	users.Use(middlewares.TokenExtractor(), authMiddleware.MiddlewareFunc())
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.UpdateMe)
	}

	// Rotas administrativas para gerenciamento de usuários
	adminUsers := router.Group("/admin/users")
	adminUsers.Use(middlewares.TokenExtractor(), authMiddleware.MiddlewareFunc(), middlewares.AdminRequired())
	{
		adminUsers.GET("/", userHandler.List)
		adminUsers.GET("/:id", userHandler.GetByID)
		adminUsers.PUT("/:id", userHandler.Update)
		adminUsers.DELETE("/:id", userHandler.Delete)
		adminUsers.PUT("/:id/promote", userHandler.PromoteToAdmin)
	}
}
