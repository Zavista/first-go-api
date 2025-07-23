package main

import (
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Balance   int64  `json:"balance"`
}

type BalanceResponse struct {
	ID      int   `json:"id"`
	Balance int64 `json:"balance"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
