package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"habit-tracker/internal/database"
	"habit-tracker/internal/models"
)

// GetHabits получает все привычки пользователя
func GetHabits(c *gin.Context) {
	// Для простоты используем фиксированный userID
	// В реальном проекте его бы брали из JWT токена
	userID := "user123"

	storage := database.GetStorage()
	habits, err := storage.GetHabitsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get habits",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"habits": habits,
	})
}

// CreateHabit создает новую привычку
func CreateHabit(c *gin.Context) {
	var req models.CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Для простоты используем фиксированный userID
	userID := "user123"

	habit := &models.Habit{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	storage := database.GetStorage()
	if err := storage.CreateHabit(habit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create habit",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"habit": habit,
	})
}

// UpdateHabit обновляет привычку
func UpdateHabit(c *gin.Context) {
	habitID := c.Param("id")
	if habitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Habit ID is required",
		})
		return
	}

	var req models.UpdateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	storage := database.GetStorage()
	if err := storage.UpdateHabit(habitID, &req); err != nil {
		if err.Error() == "habit not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Habit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update habit",
		})
		return
	}

	// Получаем обновленную привычку
	habit, _ := storage.GetHabitByID(habitID)

	c.JSON(http.StatusOK, gin.H{
		"habit": habit,
	})
}

// DeleteHabit удаляет привычку
func DeleteHabit(c *gin.Context) {
	habitID := c.Param("id")
	if habitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Habit ID is required",
		})
		return
	}

	storage := database.GetStorage()
	if err := storage.DeleteHabit(habitID); err != nil {
		if err.Error() == "habit not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Habit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete habit",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Habit deleted successfully",
	})
}

// LogHabit отмечает выполнение привычки
func LogHabit(c *gin.Context) {
	habitID := c.Param("id")
	if habitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Habit ID is required",
		})
		return
	}

	var req models.LogHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID := "user123"

	log := &models.HabitLog{
		HabitID:   habitID,
		UserID:    userID,
		Date:      req.Date,
		Completed: req.Completed,
		Notes:     req.Notes,
	}

	storage := database.GetStorage()
	if err := storage.LogHabit(log); err != nil {
		if err.Error() == "habit not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Habit not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to log habit",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"log": log,
	})
}

// GetHabitLogs получает историю выполнения привычки
func GetHabitLogs(c *gin.Context) {
	habitID := c.Param("id")
	if habitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Habit ID is required",
		})
		return
	}

	storage := database.GetStorage()
	logs, err := storage.GetHabitLogs(habitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get habit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
