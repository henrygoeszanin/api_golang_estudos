package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/henrygoeszanin/api_golang_estudos/config"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// SetupDatabase configura a conexão com o banco de dados PostgreSQL
func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao PostgreSQL: %w", err)
	}

	// Auto Migrate - cria tabelas baseadas nas entidades
	err = db.AutoMigrate(
		&entities.User{},
		&entities.Book{},
		&entities.Loan{},
	)
	if err != nil {
		return nil, fmt.Errorf("falha na migração do banco: %w", err)
	}

	return db, nil
}
