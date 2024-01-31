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
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	ec := echo.New()

	// Add static files support
	ec.Static("/static", "assets")

	// Dependencies
	dbInstance, err := infra.OpenDB()
	if err != nil {
		fmt.Println("Error opening DB")
		panic(err)
	}

	sessionStore, err := infra.NewSessionStore(dsn)
	if err != nil {
		fmt.Println("Error creating session store")
		panic(err)
	}
	defer sessionStore.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer sessionStore.StopCleanup(sessionStore.Cleanup(time.Minute * 5))

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
	csrfMiddleware := middlewares.CSRFMiddleware{}
	flashMiddleware := middlewares.FlashMiddleware{
		FlashStore: flashStore,
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, flashStore)
	demoHandler := handlers.NewDemoHandler(dbInstance)

	// Middlewares
	ec.Use(appContextMiddleware.WithAppContextMiddleware)
	ec.Use(authMiddleware.WithUserMiddleware)
	ec.Use(csrfMiddleware.WithCSRFMiddleware)
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

	ec.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	// Demo routes
	ec.GET("/", demoHandler.HandleProductIndex, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.GET("/cart", demoHandler.HandleGetCart, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.PUT("/cart", demoHandler.HandleAddToCart, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.PATCH("/cart", demoHandler.HandleDecreaseQuantity, authMiddleware.WithAuthenticationRequiredMiddleware)
	ec.POST("/cart/delete", demoHandler.HandleRemoveFromCart, authMiddleware.WithAuthenticationRequiredMiddleware)
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

	ec.Start(":8080")
	fmt.Println("Server running on port 8080")

}
