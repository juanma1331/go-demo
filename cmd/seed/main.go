package main

import (
	"context"
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/domain"
	"go-demo/internal/infra"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const NUMBER_OF_PRODUCTS = 3

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	fmt.Println("Starting seeding...")

	db, err := infra.OpenDB()
	if err != nil {
		panic(err)
	}

	migrationStart := time.Now()

	products := createProducts()

	_, insertProductsErr := db.NewInsert().Model(&products).Exec(context.Background())
	if insertProductsErr != nil {
		panic(insertProductsErr)
	}

	admin := createAdmin()

	_, insertAdminErr := db.NewInsert().Model(admin).Exec(context.Background())
	if insertAdminErr != nil {
		panic(insertAdminErr)
	}

	migrationEnd := time.Since(migrationStart)
	fmt.Printf("Seeding finished in %s\n", migrationEnd)
}

func createProducts() []domain.Product {
	images := []string{
		"./cmd/seed/images/1.jpg",
		"./cmd/seed/images/2.jpg",
		"./cmd/seed/images/3.jpg",
		"./cmd/seed/images/4.jpg",
	}

	var products []domain.Product

	for i := 0; i < NUMBER_OF_PRODUCTS; i++ {
		image := images[i%len(images)]
		img, err := imageToBytes(image)
		if err != nil {
			panic(err)
		}

		smallImage, err := app.ResizeImage(img, 200, 200)
		if err != nil {
			panic(fmt.Errorf("createProduct: failed to resize small image: %w", err))
		}

		mediumImage, err := app.ResizeImage(img, 384, 192)
		if err != nil {
			panic(fmt.Errorf("createProduct: failed to resize medium image: %w", err))
		}

		products = append(products, domain.Product{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("Product %d", i+1),
			Description: fmt.Sprintf("Product %d description", i+1),
			ImageSmall:  smallImage,
			ImageMedium: mediumImage,
		})
	}

	return products
}

func createAdmin() *domain.User {
	password, err := infra.NewBCryptPasswordManager().GenerateFromPassword("147")
	if err != nil {
		panic(err)
	}

	return &domain.User{
		ID:       uuid.New(),
		Email:    "john-doe@mail.com",
		Password: string(password),
		IsAdmin:  true,
	}
}

func imageToBytes(filepath string) ([]byte, error) {
	img, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return img, nil
}
