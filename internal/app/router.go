package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/internal/handler"
	"github.com/rubenkristian/backend/internal/middleware"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/utils"
	"gorm.io/gorm"
)

func InitializeRoute(app *fiber.App, db *gorm.DB, emailer *utils.Emailer, env *configs.EnvConfig) {
	jwtKey := env.LoadJwtConfig()
	authGenerator := utils.InitializeAuth(jwtKey.SecretKey, jwtKey.RefreshKey)
	userService := services.InitializeUserService(authGenerator, db, emailer)
	productService := services.InitializeProductService(db)

	oauth2Config := env.LoadOAuthConfig()

	authMiddleware := middleware.InitializeAuthMiddleware(userService, authGenerator)

	authHandler := handler.InitializeAuthHandler(userService, authGenerator, oauth2Config)
	productHandler := handler.InitializeProductHandler(productService)
	userHandler := handler.InitializeUserHandler(userService)

	api := app.Group("/api", authMiddleware.CheckAuthorization)

	api.Get("/users/:user_id", authMiddleware.CheckRole([]string{"user", "admin"}), userHandler.GetUser)
	api.Get("/users", authMiddleware.CheckRole([]string{"user", "admin"}), userHandler.GetAllUser)
	api.Post("/users", authMiddleware.CheckRole([]string{"admin"}), userHandler.CreateUser)
	api.Put("/users/:user_id", authMiddleware.CheckRole([]string{"admin"}), userHandler.UpdateUser)
	api.Delete("/users/:user_id", authMiddleware.CheckRole([]string{"admin"}), userHandler.DeleteUser)

	api.Get("/products/:product_id", authMiddleware.CheckRole([]string{"user", "admin"}), productHandler.GetProduct)
	api.Get("/products", authMiddleware.CheckRole([]string{"user", "admin"}), productHandler.GetAllProduct)
	api.Post("/products", authMiddleware.CheckRole([]string{"admin"}), productHandler.PostCreateProduct)
	api.Put("/products/:product_id", authMiddleware.CheckRole([]string{"admin"}), productHandler.UpdateProduct)
	api.Delete("/products/:product_id", authMiddleware.CheckRole([]string{"admin"}), productHandler.DeleteProduct)

	auth := app.Group("/auth")

	auth.Post("/user/login", authHandler.LoginUser)
	auth.Post("/admin/login", authHandler.LoginAdmin)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/google/login", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	app.Post("/register", authHandler.RegisterUser)
}