package main

import (
	"log"
	"os"
	"strings"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/handler"
	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := gin.Default()

	// CORS origins: comma-separated from CORS_ORIGINS env, fallback to localhost
	allowOrigins := []string{"http://localhost:5173", "http://localhost:4173"}
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		allowOrigins = strings.Split(origins, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		// Public auth routes
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)

		// Protected routes (require auth)
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", handler.Me)

			protected.GET("/categories", handler.GetCategories)
			protected.GET("/categories/:id", handler.GetCategory)
			protected.POST("/categories", handler.CreateCategory)
			protected.PUT("/categories/:id", handler.UpdateCategory)
			protected.DELETE("/categories/:id", handler.DeleteCategory)

			protected.GET("/transactions", handler.GetTransactions)
			protected.GET("/transactions/:id", handler.GetTransaction)
			protected.POST("/transactions", handler.CreateTransaction)
			protected.PUT("/transactions/:id", handler.UpdateTransaction)
			protected.DELETE("/transactions/:id", handler.DeleteTransaction)

			protected.GET("/reports/summary", handler.GetSummary)
			protected.GET("/reports/yearly-trend", handler.GetYearlyTrend)
			protected.GET("/reports/category-trend", handler.GetCategoryTrend)

			// Savings goals
			protected.GET("/savings-goals", handler.GetSavingsGoals)
			protected.GET("/savings-goals/:id", handler.GetSavingsGoal)
			protected.POST("/savings-goals", handler.CreateSavingsGoal)
			protected.PUT("/savings-goals/:id", handler.UpdateSavingsGoal)
			protected.DELETE("/savings-goals/:id", handler.DeleteSavingsGoal)
			protected.POST("/savings-goals/:id/deposit", handler.DepositToSavingsGoal)
			protected.POST("/savings-goals/:id/withdraw", handler.WithdrawFromSavingsGoal)

			// Budgets
			protected.GET("/budgets", handler.GetBudgets)
			protected.GET("/budgets/summary", handler.GetBudgetSummary)
			protected.POST("/budgets", handler.UpsertBudget)
			protected.DELETE("/budgets/:id", handler.DeleteBudget)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
