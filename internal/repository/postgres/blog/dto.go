package blog

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
	Author  string `json:"author"`
}

type AdminGetDetail struct {
	bun.BaseModel `bun:"table:blogs"`

	Id      string `json:"id" bun:"id"`
	Title   string `json:"title" bun:"title"`
	Content string `json:"content" bun:"content"`
	Author  string `json:"author" bun:"author"`
}

type AdminCreateRequest struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Author  string `json:"author" form:"author"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:blogs"`
	Id            string    `json:"id" bun:"id"`
	Title         string    `json:"title" bun:"title"`
	Content       string    `json:"content" bun:"content"`
	Author        string    `json:"author" bun:"author"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     *string   `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	Id      string  `json:"id" form:"id"`
	Title   *string `json:"title" bun:"title"`
	Content *string `json:"content" bun:"content"`
	Author  *string `json:"author" bun:"author"`
}

type AdminDetailResponseSwagger struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type AdminCreateResponseSwagger struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}
