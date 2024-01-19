package main

import (
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/app/handlers"
	"go-demo/internal/app/middlewares"
	"go-demo/internal/app/services"
	"go-demo/internal/app/services/authservice"
	"go-demo/internal/infra"
	"go-demo/internal/infra/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	ec := echo.New()

	// Add static files support
	ec.Static("/static", "assets")

	// Dependencies
	dbInstance, err := infra.OpenDB(infra.DSN)
	if err != nil {
		panic(err)
	}

	sessionStore, err := infra.NewBunSessionStore(dbInstance, []byte("session-secret-key"))
	if err != nil {
		panic(err)
	}

	flashStore := app.NewFlashStore("flash-secret-key")

	authTokenManager := infra.NewAuthTokenManager()
	passwordManager := infra.NewBCryptPasswordManager()

	// Repositories
	userRepository := repositories.NewSqliteUserRepository(dbInstance)
	authTokenRepository := repositories.NewSqliteAuthTokenRepository(dbInstance)

	// Validation
	uniqueEmailValidation := authservice.UniqueEmailValidator{
		UserRepository: userRepository,
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	validator := services.NewPlaygroundValidator(v)

	validator.RegisterValidation(
		authservice.ValidatorUniqueEmailKey,
		uniqueEmailValidation.UniqueEmailValidation,
		authservice.ValidatorUniqueEmailErrorMsg,
	)

	// Services
	authService := authservice.NewAuthService(authservice.AuthServiceParams{
		UserRepository:   userRepository,
		AuthTokenRepo:    authTokenRepository,
		PasswordManager:  passwordManager,
		AuthTokenManager: authTokenManager,
		SessionStore:     sessionStore,
		Validator:        validator,
	})

	// Middlewares
	appContextMiddleware := middlewares.AppContextMiddleware{}
	authMiddleware := middlewares.AuthMiddleware{
		SessionStore:     sessionStore,
		UserRepo:         userRepository,
		AuthTokenRepo:    authTokenRepository,
		AuthTokenManager: authTokenManager,
	}
	flashMiddleware := middlewares.FlashMiddleware{
		FlashStore: flashStore,
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, flashStore)
	productHandler := handlers.NewProductHandler()

	// Middlewares
	ec.Use(appContextMiddleware.WithAppContextMiddleware)
	ec.Use(authMiddleware.WithUserMiddleware)
	ec.Use(flashMiddleware.WithFlashMiddleWare)

	ec.Use(func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(app.AppContext)
			fmt.Printf("User: %+v\n", cc.User)
			fmt.Printf("Flash: %+v\n", cc.Flash)
			return hf(cc)
		}
	})

	ec.Use(middleware.Gzip())

	// Product routes
	ec.GET("/", productHandler.HandleProductIndex)
	// app.GET("/products", productHandler.HandleProductIndex)
	// app.GET("/products/:id", productHandler.HandleProductDetail)
	// app.GET("/products/:id/image", productHandler.HandleProductImage)

	ec.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	// Auth routes
	ag := ec.Group("/auth")
	ag.GET("/register", authHandler.HandleShowRegister)
	ag.GET("/login", authHandler.HandleShowLogin)
	ag.POST("/register", authHandler.HandleRegister)
	ag.POST("/login", authHandler.HandleLogin)
	ag.GET("/logout", authHandler.HandleLogout)

	ec.Start(":8080")
	fmt.Println("Server running on port 8080")

}
