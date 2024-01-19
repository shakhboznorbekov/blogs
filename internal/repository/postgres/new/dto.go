package new

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	Title  *string
}

type AdminGetListResponse struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Source  string `json:"source"`
}

type AdminGetDetail struct {
	bun.BaseModel `bun:"table:news"`

	Id      string `json:"id" bun:"id"`
	Title   string `json:"title" bun:"title"`
	Content string `json:"content" bun:"content"`
	Source  string `json:"source" bun:"source"`
}

type AdminCreateRequest struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Source  string `json:"source" form:"source"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:users"`
	Id            string    `json:"id" bun:"id"`
	Title         string    `json:"title" bun:"title"`
	Content       string    `json:"content" bun:"content"`
	Source        string    `json:"source" bun:"source"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     *string   `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	Id      string  `json:"id" form:"id"`
	Title   *string `json:"title" form:"title"`
	Content *string `json:"content" form:"content"`
	Source  *string `json:"source" form:"source"`
}

type AdminDetailResponseSwagger struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Source  string `json:"source"`
}

type AdminCreateResponseSwagger struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Source  string `json:"source"`
}
