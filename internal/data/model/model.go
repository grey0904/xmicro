package model

import "time"

type Order struct {
	ID          int
	UserID      int
	CustomerID  int
	TotalAmount float64
	OrderDate   time.Time
	Status      string
}

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}
