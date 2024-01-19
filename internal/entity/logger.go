package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Logger struct {
	bun.BaseModel `bun:"table:loggers"`
	Id            string      `json:"id" bun:"id,pk,autoincrement"`
	Action        string      `json:"action" bun:"action"`
	Data          interface{} `json:"data"   bun:"data"`
	Method        string      `json:"method" bun:"method"`
	CreatedAt     time.Time   `json:"-" bun:"created_at"`
}

type LogCreateDto struct {
	Action    string      `json:"action" form:"action"`
	Data      interface{} `json:"data" form:"data"`
	Method    string      `json:"method" form:"method"`
	CreatedAt time.Time   `json:"-" bun:"created_at"`
}
