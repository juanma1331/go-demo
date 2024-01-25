package main

import (
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/app/handlers"
	"go-demo/internal/app/middlewares"
	"go-demo/internal/app/services"
	"go-demo/internal/app/services/authservice"
	auth_validation "go-demo/internal/app/services/authservice/validation"
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

	passwordManager := infra.NewBCryptPasswordManager()

	// Repositories
	userRepository := repositories.NewSqliteUserRepository(dbInstance)

	// Validation
	v := validator.New(validator.WithRequiredStructEnabled())
	validator := services.NewPlaygroundValidator(v)
	auth_validation.RegisterUniqueEmailValidator(&validator, userRepository)
	auth_validation.RegisterPasswordValidations(&validator)

	// Services
	authService := authservice.NewAuthService(authservice.AuthServiceParams{
		UserRepository:  userRepository,
		PasswordManager: passwordManager,
		SessionStore:    sessionStore,
		Validator:       &validator,
	})

	// Middlewares
	appContextMiddleware := middlewares.AppContextMiddleware{}
	authMiddleware := middlewares.AuthMiddleware{
		SessionStore: sessionStore,
		UserRepo:     userRepository,
	}
	flashMiddleware := middlewares.FlashMiddleware{
		FlashStore: flashStore,
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, flashStore)
	demoHandler := handlers.NewDemoHandler(dbInstance)

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
	ec.Use(middleware.Recover())

	// Demo routes
	ec.GET("/", demoHandler.HandleProductIndex, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.GET("/cart", demoHandler.HandleGetCart, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.PUT("/cart", demoHandler.HandleAddToCart, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.GET("/products", demoHandler.GetProductList, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.GET("/products/:id/image/:size", demoHandler.HandleProductImage, authMiddleware.WithAuthenticationRequiredMiddleware)

	// Auth routes
	ag := ec.Group("/auth")
	ag.GET("/register", authHandler.HandleShowRegister)
	ag.POST("/register", authHandler.HandleRegister)
	ag.GET("/register/validate-email", authHandler.HandleValidateRegisterEmail)
	ag.GET("/register/validate-password", authHandler.HandleValidateRegisterPassword)

	ag.GET("/login", authHandler.HandleShowLogin)
	ag.POST("/login", authHandler.HandleLogin)
	ag.GET("/logout", authHandler.HandleLogout)

	ec.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	ec.Start(":8080")
	fmt.Println("Server running on port 8080")

}
