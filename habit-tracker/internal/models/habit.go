package models

import (
	"time"
)

// Habit представляет привычку пользователя
type Habit struct {
	ID          string    `json:"id" dynamodbav:"id"`
	UserID      string    `json:"user_id" dynamodbav:"user_id"`
	Name        string    `json:"name" dynamodbav:"name"`
	Description string    `json:"description,omitempty" dynamodbav:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

// HabitLog представляет запись о выполнении привычки
type HabitLog struct {
	HabitID   string    `json:"habit_id" dynamodbav:"habit_id"`
	UserID    string    `json:"user_id" dynamodbav:"user_id"`
	Date      string    `json:"date" dynamodbav:"date"` // format: 2025-01-15
	Completed bool      `json:"completed" dynamodbav:"completed"`
	Notes     string    `json:"notes,omitempty" dynamodbav:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"created_at"`
}

// CreateHabitRequest для создания новой привычки
type CreateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// LogHabitRequest для логирования выполнения привычки
type LogHabitRequest struct {
	Date      string `json:"date" binding:"required"` // format: 2025-01-15
	Completed bool   `json:"completed"`
	Notes     string `json:"notes"`
}

// UpdateHabitRequest для обновления привычки
type UpdateHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
