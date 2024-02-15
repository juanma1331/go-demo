package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	auth_handlers "github.com/juanma1331/go-demo/internal/auth/app/handlers"
	auth_middlewares "github.com/juanma1331/go-demo/internal/auth/app/middlewares"
	auth_service "github.com/juanma1331/go-demo/internal/auth/app/services"
	auth_validation "github.com/juanma1331/go-demo/internal/auth/app/services/validation"
	auth_infra "github.com/juanma1331/go-demo/internal/auth/infra"
	ecommerce_handlers "github.com/juanma1331/go-demo/internal/ecommerce/app/handlers"

	"os"
	"time"

	"github.com/juanma1331/go-demo/internal/shared"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//go:embed assets/*
var embeddedFiles embed.FS

func main() {
	assetsFS, err := fs.Sub(embeddedFiles, "assets")
	if err != nil {
		log.Fatal("failed to locate embedded assets:", err)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	ec := echo.New()

	// Dependencies
	dbInstance, err := shared.OpenDB()
	if err != nil {
		fmt.Println("Error opening DB")
		panic(err)
	}
	defer dbInstance.Close()

	sessionStore, err := auth_infra.NewSessionStore(dsn)
	if err != nil {
		fmt.Println("Error creating session store")
		panic(err)
	}
	defer sessionStore.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer sessionStore.StopCleanup(sessionStore.Cleanup(time.Minute * 5))

	flashStore := shared.NewFlashStore("flash-secret-key")

	passwordManager := auth_infra.NewBCryptPasswordManager()

	// Repositories
	userRepository := auth_infra.NewUserRepository(dbInstance)

	// Validation
	v := validator.New(validator.WithRequiredStructEnabled())
	validator := shared.NewPlaygroundValidator(v)
	auth_validation.RegisterUniqueEmailValidator(&validator, userRepository)
	auth_validation.RegisterPasswordValidations(&validator)

	// Services
	authService := auth_service.NewAuthService(auth_service.AuthServiceParams{
		UserRepository:  userRepository,
		PasswordManager: passwordManager,
		SessionStore:    sessionStore,
		Validator:       &validator,
	})

	// Middlewares
	appContextMiddleware := shared.AppContextMiddleware{}
	authMiddleware := auth_middlewares.AuthMiddleware{
		SessionStore: sessionStore,
		UserRepo:     userRepository,
	}
	csrfMiddleware := shared.CSRFMiddleware{}
	flashMiddleware := shared.FlashMiddleware{
		FlashStore: flashStore,
	}

	// Handlers
	showLoginHandler := auth_handlers.NewShowLoginHandler()
	loginHandler := auth_handlers.NewLoginHandler(authService, flashStore)
	showRegisterHandler := auth_handlers.NewShowRegisterHandler()
	registerHandler := auth_handlers.NewRegisterHandler(authService, flashStore)
	validateRegisterEmailHandler := auth_handlers.NewValidateRegisterEmailHandler(authService)
	validateRegisterPassword := auth_handlers.NewValidateRegisterPasswordHandler(authService)
	logoutHandler := auth_handlers.NewLogoutHandler(authService, flashStore)

	showProductIndexHandler := ecommerce_handlers.NewShowProductIndexHandler()
	getProductListHandler := ecommerce_handlers.NewGetProductListHandler(dbInstance)
	getMoreProductsHandler := ecommerce_handlers.NewGetMoreProductsHandler(dbInstance)
	getCartHandler := ecommerce_handlers.NewGetCartHandler(dbInstance)
	addToCartHandler := ecommerce_handlers.NewAddToCartHandler(dbInstance)
	decreaseQuantityHandler := ecommerce_handlers.NewDecreaseQuantityHandler(dbInstance)
	removeFromCartHandler := ecommerce_handlers.NewRemoveFromCartHandler(dbInstance)
	getProductImageHandler := ecommerce_handlers.NewGetProductImageHandler(dbInstance)

	// Middlewares
	ec.Use(appContextMiddleware.WithAppContextMiddleware)
	ec.Use(authMiddleware.WithUserMiddleware)
	ec.Use(csrfMiddleware.WithCSRFMiddleware)
	ec.Use(flashMiddleware.WithFlashMiddleWare)

	// Error propagation
	ec.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err != nil {
				c.Error(err)
			}
			return nil
		}
	})
	ec.Use(middleware.Gzip())
	ec.Use(middleware.Recover())

	ec.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(assetsFS)))))

	// Routes
	ec.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(302, "/products")
	})

	ec.HTTPErrorHandler = shared.CustomHTTPErrorHandler
	// Ecommerce routes
	cartGroup := ec.Group("/cart")
	cartGroup.GET("", getCartHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)
	cartGroup.PUT("", addToCartHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)
	cartGroup.PATCH("", decreaseQuantityHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)
	// This should be a DELETE request, but we are using POST due to problems with csrf
	cartGroup.POST("/delete", removeFromCartHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)

	productsGroup := ec.Group("/products")
	productsGroup.GET("", showProductIndexHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)
	productsGroup.GET("/get-list", getProductListHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)
	productsGroup.GET("/get-more/:cursor", getMoreProductsHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)
	productsGroup.GET("/:id/image/:size", getProductImageHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)

	// Auth routes
	ag := ec.Group("/auth")
	ag.GET("/register", showRegisterHandler.Handler)
	ag.POST("/register", registerHandler.Handler)
	ag.GET("/register/validate-email", validateRegisterEmailHandler.Handler)
	ag.GET("/register/validate-password", validateRegisterPassword.Handler)

	ag.GET("/login", showLoginHandler.Handler)
	ag.POST("/login", loginHandler.Handler)
	ag.GET("/logout", logoutHandler.Handler, authMiddleware.WithAuthenticationRequiredMiddleware)

	ec.Start(":8080")
	fmt.Println("Server running on port 8080")

}
