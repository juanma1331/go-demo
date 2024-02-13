package main

import (
	"context"
	"fmt"
	auth_domain "go-demo/internal/auth/domain"
	auth_infra "go-demo/internal/auth/infra"
	ecommerce_domain "go-demo/internal/ecommerce/domain"
	"go-demo/internal/shared"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const NUMBER_OF_PRODUCTS = 1000

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	fmt.Println("Starting seeding...")

	db, err := shared.OpenDB()
	if err != nil {
		panic(err)
	}

	seedingStart := time.Now()

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

	seedingEnd := time.Since(seedingStart)
	fmt.Printf("Seeding finished in %s\n", seedingEnd)
}

func createProducts() []ecommerce_domain.Product {
	images := []string{
		"./cmd/seed/images/1.jpg",
		"./cmd/seed/images/2.jpg",
		"./cmd/seed/images/3.jpg",
		"./cmd/seed/images/4.jpg",
	}

	var products []ecommerce_domain.Product

	for i := 0; i < NUMBER_OF_PRODUCTS; i++ {
		image := images[i%len(images)]
		img, err := imageToBytes(image)
		if err != nil {
			panic(err)
		}

		smallImage, err := shared.ResizeImage(img, 64, 64)
		if err != nil {
			panic(fmt.Errorf("createProduct: failed to resize small image: %w", err))
		}

		mediumImage, err := shared.ResizeImage(img, 180, 192)
		if err != nil {
			panic(fmt.Errorf("createProduct: failed to resize medium image: %w", err))
		}

		products = append(products, ecommerce_domain.Product{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("Product %d", i+1),
			Description: fmt.Sprintf("Product %d description", i+1),
			Price:       int64(rand.Intn(1000)),
			ImageSmall:  smallImage,
			ImageMedium: mediumImage,
		})
	}

	return products
}

func createAdmin() *auth_domain.User {
	password, err := auth_infra.NewBCryptPasswordManager().GenerateFromPassword("147")
	if err != nil {
		panic(err)
	}

	return &auth_domain.User{
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
