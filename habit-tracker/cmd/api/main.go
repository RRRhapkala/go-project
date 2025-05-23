package main

import (
	"log"
	"net/http"
	"os"

	"habit-tracker/internal/config"
	"habit-tracker/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем роутер
	router := gin.Default()

	// Middleware для CORS (для фронтенда)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Habits endpoints
		api.GET("/habits", handlers.GetHabits)
		api.POST("/habits", handlers.CreateHabit)
		api.PUT("/habits/:id", handlers.UpdateHabit)
		api.DELETE("/habits/:id", handlers.DeleteHabit)

		// Habit logs endpoints
		api.POST("/habits/:id/log", handlers.LogHabit)
		api.GET("/habits/:id/logs", handlers.GetHabitLogs)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "habit-tracker",
		})
	})

	// Простая главная страница (пока без шаблонов)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Habit Tracker API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "/health",
				"habits": "/api/v1/habits",
			},
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
