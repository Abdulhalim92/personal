package models

import "time"

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// User пользователи
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Token токены
type Token struct {
	ID     int
	UserID int
	Token  string
}

// Account счеты пользователя
type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserId    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Type типы операций
type Type struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Category категории операций
type Category struct {
	ID        int       `json:"id"`
	TypeId    int       `json:"type_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Operation типы операций
type Operation struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TypeID      int       `json:"type_id"`
	CategoryID  int       `json:"category_id"`
	AccountID   int       `json:"account_id"`
	AccountIDTo int       `json:"account_to"` // Для операции трансфера
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

// Report получение отчета
type Report struct {
	AccountID int    `json:"account_id"`
	TypeID    int    `json:"type_id"`
	Page      int    `json:"page"`
	Limit     int    `json:"count_in_page"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// EditDelAccount редактирование и удаление счетов
type EditDelAccount struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
}

type GetReports struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	AccountID  int       `json:"account_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
