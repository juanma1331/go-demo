package main

import (
	"context"
	"go-demo/internal/domain"
	"go-demo/internal/infra"

	"github.com/google/uuid"
)

func main() {
	db, err := infra.OpenDB(infra.DSN)
	if err != nil {
		panic(err)
	}

	// products := createProducts()

	// _, insertProductsErr := db.NewInsert().Model(&products).Exec(context.Background())
	// if insertProductsErr != nil {
	// 	panic(insertProductsErr)
	// }

	admin := createAdmin()

	_, insertAdminErr := db.NewInsert().Model(admin).Exec(context.Background())
	if insertAdminErr != nil {
		panic(insertAdminErr)
	}
}

// func createProducts() []model.Product {
// 	images := []string{
// 		"./cmd/seed/images/1.jpg",
// 		"./cmd/seed/images/2.jpg",
// 		"./cmd/seed/images/3.jpg",
// 		"./cmd/seed/images/4.jpg",
// 	}

// 	var products []model.Product

// 	for i, image := range images {
// 		img, err := imageToBytes(image)
// 		if err != nil {
// 			panic(err)
// 		}

// 		products = append(products, model.Product{
// 			ID:          uuid.New(),
// 			Name:        fmt.Sprintf("Product %d", i+1),
// 			Description: fmt.Sprintf("Product %d description", i+1),
// 			Image:       img,
// 		})
// 	}

// 	return products
// }

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

// func imageToBytes(filepath string) ([]byte, error) {
// 	img, err := os.ReadFile(filepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return img, nil
// }
