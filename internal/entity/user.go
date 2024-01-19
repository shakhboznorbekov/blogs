package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id        string     `json:"id" bun:"id"`
	Username  string     `json:"username" bun:"username"`
	Password  *string    `json:"password" bun:"password"`
	Role      string     `json:"role" bun:"role"`
	Status    bool       `json:"status" bun:"status"`
	CreatedAt time.Time  `json:"-" bun:"created_at"`
	CreatedBy *string    `json:"-" bun:"created_by"`
	UpdatedAt *time.Time `json:"-" bun:"updated_at"`
	UpdatedBy *string    `json:"-" bun:"updated_by"`
	DeletedAt *time.Time `json:"-" bun:"deleted_at"`
	DeletedBy *string    `json:"-" bun:"deleted_by"`
}

// Role
//1. Admin
