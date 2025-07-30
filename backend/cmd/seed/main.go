package main

import (
	"context"
	"log"
	"server/internal/repository"
	"server/internal/service/seed"
)

func main() {
	ctx := context.Background()
	dbRepo, err := repository.NewDbRepo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	seedService := seed.NewSeedService(dbRepo)
	if err := seedService.Seed(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Seed completed successfully")
}
