package user

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Username *string
}

type AdminGetListResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Status   *bool  `json:"status"`
	Role     string `json:"role"`
}

type AdminGetDetail struct {
	bun.BaseModel `bun:"table:users"`

	Id       string  `json:"id" bun:"id"`
	Username string  `json:"username" bun:"username"`
	Password *string `json:"-" bun:"password"`
	Status   bool    `json:"status" bun:"status"`
	Role     string  `json:"role" bun:"role"`
}

type AdminCreateRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Status   bool   `json:"status" form:"status"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:users"`
	Id            string    `json:"id" bun:"id"`
	Username      string    `json:"username" bun:"username"`
	Password      string    `json:"password" bun:"password"`
	Role          *string   `json:"role" bun:"role"`
	Status        bool      `json:"status" bun:"status"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     *string   `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	Id       string  `json:"id" form:"id"`
	Username *string `json:"username" form:"username"`
	Role     *string `json:"role" form:"role"`
	Status   *bool   `json:"status" form:"status"`
}

type AdminDetailResponseSwagger struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	Password *string `json:"-"`
	Status   bool    `json:"status"`
	Role     string  `json:"role"`
}

type AdminCreateResponseSwagger struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Role     *string `json:"role"`
	Status   bool    `json:"status"`
}
