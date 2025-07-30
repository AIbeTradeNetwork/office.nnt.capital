package main

import (
	"context"
	"log"
	"server/internal/repository/mongodb"
	"server/internal/service/seed/data"
)

func main() {
	ctx := context.Background()

	// Подключение к MongoDB
	repo, err := mongodb.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Получаем продукты для добавления
	products := data.SubscriptionProducts()

	for _, product := range products {
		// Проверяем, существует ли уже продукт
		existingProduct, err := repo.ProductGetByCode(ctx, product.Code)
		if err == nil && existingProduct != nil {
			log.Printf("Product %s already exists, skipping", product.Code)
			continue
		}

		// Создаём продукт
		err = repo.ProductCreate(ctx, &product)
		if err != nil {
			log.Printf("Failed to create product %s: %v", product.Code, err)
		} else {
			log.Printf("Successfully created product: %s", product.Code)
		}
	}

	// Закрываем соединение
	repo.Disconnect(ctx)

	log.Println("Finished adding subscription products")
}
