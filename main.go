package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/henrygoeszanin/api_golang_estudos/config"
	"github.com/henrygoeszanin/api_golang_estudos/infrastructure/database"
	"github.com/henrygoeszanin/api_golang_estudos/presentation/routes"
)

func init() {
	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis padrão")
	}
}

func main() {
	// Carregar configurações
	cfg := config.LoadConfig()

	// Inicializar o banco de dados
	db, err := database.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Configurar o router Gin
	router := gin.Default()

	// Configurar rotas
	routes.SetupRoutes(router, db, cfg)

	// Iniciar o servidor
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Servidor iniciado na porta %s\n", cfg.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
