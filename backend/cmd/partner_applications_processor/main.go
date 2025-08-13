package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/repository"
	"server/internal/service/team"
)

func main() {
	log.Println("Starting Partner Applications Processor...")

	// Подключаемся к MongoDB
	dbCtx := context.Background()
	mongoRepo, err := repository.NewDbRepo(dbCtx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoRepo.Disconnect(dbCtx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Создаем team service
	teamService := team.NewTeamService(mongoRepo, nil, nil)

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатываем сигналы для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping...")
		cancel()
	}()

	// Запускаем обработчик просроченных заявок
	ticker := time.NewTicker(30 * time.Minute) // Каждые 30 минут
	defer ticker.Stop()

	// Выполняем первую обработку сразу
	log.Println("Running initial expired applications processing...")
	if err := teamService.ProcessExpiredApplications(ctx); err != nil {
		log.Printf("Error in initial processing: %v", err)
	}

	// Основной цикл
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, shutting down...")
			return

		case <-ticker.C:
			log.Println("Running scheduled expired applications processing...")
			if err := teamService.ProcessExpiredApplications(ctx); err != nil {
				log.Printf("Error processing expired applications: %v", err)
			}
		}
	}
}
