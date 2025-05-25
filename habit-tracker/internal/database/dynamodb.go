package database

import (
	"errors"
	"sync"
	"time"

	"habit-tracker/internal/models"
)

// Для начала используем in-memory storage, потом заменим на DynamoDB
type HabitsStorage struct {
	habits map[string]*models.Habit
	logs   map[string][]*models.HabitLog // key: habitID
	mu     sync.RWMutex
}

var storage *HabitsStorage
var once sync.Once

// GetStorage возвращает singleton экземпляр хранилища
func GetStorage() *HabitsStorage {
	once.Do(func() {
		storage = &HabitsStorage{
			habits: make(map[string]*models.Habit),
			logs:   make(map[string][]*models.HabitLog),
		}
	})
	return storage
}

// CreateHabit создает новую привычку
func (s *HabitsStorage) CreateHabit(habit *models.Habit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	habit.CreatedAt = time.Now()
	habit.UpdatedAt = time.Now()
	s.habits[habit.ID] = habit
	return nil
}

// GetHabitsByUserID получает все привычки пользователя
func (s *HabitsStorage) GetHabitsByUserID(userID string) ([]*models.Habit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userHabits []*models.Habit
	for _, habit := range s.habits {
		if habit.UserID == userID {
			userHabits = append(userHabits, habit)
		}
	}
	return userHabits, nil
}

// GetHabitByID получает привычку по ID
func (s *HabitsStorage) GetHabitByID(habitID string) (*models.Habit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	habit, exists := s.habits[habitID]
	if !exists {
		return nil, errors.New("habit not found")
	}
	return habit, nil
}

// UpdateHabit обновляет привычку
func (s *HabitsStorage) UpdateHabit(habitID string, updates *models.UpdateHabitRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	habit, exists := s.habits[habitID]
	if !exists {
		return errors.New("habit not found")
	}

	if updates.Name != "" {
		habit.Name = updates.Name
	}
	if updates.Description != "" {
		habit.Description = updates.Description
	}
	habit.UpdatedAt = time.Now()

	s.habits[habitID] = habit
	return nil
}

// DeleteHabit удаляет привычку
func (s *HabitsStorage) DeleteHabit(habitID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.habits[habitID]; !exists {
		return errors.New("habit not found")
	}

	delete(s.habits, habitID)
	delete(s.logs, habitID) // Удаляем и все логи
	return nil
}

// LogHabit создает запись о выполнении привычки
func (s *HabitsStorage) LogHabit(log *models.HabitLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли привычка
	if _, exists := s.habits[log.HabitID]; !exists {
		return errors.New("habit not found")
	}

	log.CreatedAt = time.Now()

	// Проверяем, есть ли уже запись на эту дату
	logs := s.logs[log.HabitID]
	for i, existingLog := range logs {
		if existingLog.Date == log.Date {
			// Обновляем существующую запись
			logs[i] = log
			s.logs[log.HabitID] = logs
			return nil
		}
	}

	// Добавляем новую запись
	s.logs[log.HabitID] = append(s.logs[log.HabitID], log)
	return nil
}

// GetHabitLogs получает все записи для привычки
func (s *HabitsStorage) GetHabitLogs(habitID string) ([]*models.HabitLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logs, exists := s.logs[habitID]
	if !exists {
		return []*models.HabitLog{}, nil
	}
	return logs, nil
}
