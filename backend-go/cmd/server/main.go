package main

// @title Agroboard Content API
// @version 1.0
// @description API backend em Go para recomendações de conteúdo.
// @host localhost:8080
// @BasePath /api
// @schemes http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"backend-go/config"
	"backend-go/database"
	_ "backend-go/docs"
	"backend-go/handler"
	"backend-go/repository"
	"backend-go/routes"
	"backend-go/service"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente (opcional, mas útil em dev)
	_ = godotenv.Load(".env")

	// Configurações da aplicação
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("erro carregando config: %v", err)
	}

	// Conexão com o banco
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	// Migração e seeds
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("erro ao migrar: %v", err)
	}
	if err := database.Seed(db); err != nil {
		log.Fatalf("erro ao semear dados: %v", err)
	}

	// Injeção de dependências
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := routes.NewRouter(userHandler).SetupRoutes()

	// Sobe o servidor em goroutine para permitir shutdown graceful
	go func() {
		addr := fmt.Sprintf(":%d", cfg.HTTPPort)
		log.Printf("servidor ouvindo em %s", addr)
		if err := router.Run(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	// Espera sinais de encerramento (Ctrl+C, kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("desligando servidor...")

	// Fecha conexões do GORM
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
