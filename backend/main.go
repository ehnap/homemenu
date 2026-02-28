package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/homemenu/backend/config"
	"github.com/homemenu/backend/db"
	"github.com/homemenu/backend/handler"
	"github.com/homemenu/backend/middleware"
	"github.com/homemenu/backend/repository/sqlite"
	"github.com/homemenu/backend/service"
	"github.com/homemenu/backend/service/generator"
	"github.com/homemenu/backend/service/llm"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.InitSQLite(cfg.DB.Path)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	// Repositories
	userRepo := sqlite.NewUserRepo(database)
	recipeRepo := sqlite.NewRecipeRepo(database)
	ingredientRepo := sqlite.NewIngredientRepo(database)
	mealPlanRepo := sqlite.NewMealPlanRepo(database)
	mealPlanItemRepo := sqlite.NewMealPlanItemRepo(database)
	settingsRepo := sqlite.NewSettingsRepo(database)

	// Generator
	var gen generator.MenuGenerator
	if cfg.LLM.APIKey != "" && cfg.LLM.BaseURL != "" {
		gen = generator.NewAIEngine(cfg.LLM.BaseURL, cfg.LLM.APIKey, cfg.LLM.Model)
	} else {
		gen = generator.NewRuleEngine()
	}

	// Services
	authService := service.NewAuthService(userRepo, cfg.Server.JWTSecret)
	recipeService := service.NewRecipeService(recipeRepo, ingredientRepo)
	mealPlanService := service.NewMealPlanService(mealPlanRepo, mealPlanItemRepo, recipeRepo, ingredientRepo, gen)
	shoppingService := service.NewShoppingService(mealPlanItemRepo, recipeRepo, ingredientRepo)
	settingsService := service.NewSettingsService(settingsRepo)
	llmClient := llm.NewClient()
	parseService := service.NewParseService(llmClient, settingsService, cfg.LLM)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	recipeHandler := handler.NewRecipeHandler(recipeService, parseService)
	mealPlanHandler := handler.NewMealPlanHandler(mealPlanService, shoppingService)
	shareHandler := handler.NewShareHandler(mealPlanService, shoppingService, recipeService)
	uploadHandler := handler.NewUploadHandler(cfg.Upload.Dir, cfg.Upload.MaxSizeMB)
	settingsHandler := handler.NewSettingsHandler(settingsService)

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve uploaded files
	r.Static("/uploads", cfg.Upload.Dir)

	// API routes
	api := r.Group("/api")
	{
		// Health check (no auth)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// Share (no auth required)
		api.GET("/share/:token", shareHandler.GetByToken)
		api.GET("/share/recipe/:token", shareHandler.GetRecipeByToken)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.Server.JWTSecret))
		{
			recipes := protected.Group("/recipes")
			{
				recipes.GET("", recipeHandler.List)
				recipes.POST("", recipeHandler.Create)
				recipes.POST("/parse-text", recipeHandler.ParseText)
				recipes.GET("/:id", recipeHandler.GetByID)
				recipes.PUT("/:id", recipeHandler.Update)
				recipes.DELETE("/:id", recipeHandler.Delete)
			recipes.POST("/:id/share", recipeHandler.GenerateShareToken)
		}

			protected.GET("/ingredients/suggestions", recipeHandler.SuggestIngredients)

			mealPlans := protected.Group("/meal-plans")
			{
				mealPlans.GET("", mealPlanHandler.List)
				mealPlans.POST("", mealPlanHandler.Create)
				mealPlans.GET("/:id", mealPlanHandler.GetByID)
				mealPlans.PUT("/:id", mealPlanHandler.Update)
				mealPlans.DELETE("/:id", mealPlanHandler.Delete)
				mealPlans.POST("/generate", mealPlanHandler.Generate)
				mealPlans.PUT("/:id/items", mealPlanHandler.UpdateItems)
				mealPlans.POST("/:id/items/:itemId/reroll", mealPlanHandler.RerollItem)
				mealPlans.GET("/:id/shopping-list", mealPlanHandler.ShoppingList)
			}

			protected.POST("/upload", uploadHandler.Upload)

			settings := protected.Group("/settings")
			{
				settings.GET("/llm", settingsHandler.GetLLMSettings)
				settings.PUT("/llm", settingsHandler.UpdateLLMSettings)
			}
		}
	}

	// Serve frontend static files
	staticSub, err := fs.Sub(staticFS, "static")
	if err == nil {
		indexHTML, _ := fs.ReadFile(staticSub, "index.html")
		r.NoRoute(func(c *gin.Context) {
			// Try to serve static file
			f, err := http.FS(staticSub).Open(c.Request.URL.Path)
			if err == nil {
				f.Close()
				http.FileServer(http.FS(staticSub)).ServeHTTP(c.Writer, c.Request)
				return
			}
			// SPA fallback: serve index.html directly to avoid http.FileServer redirect
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
		})
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("HomeMenu server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
