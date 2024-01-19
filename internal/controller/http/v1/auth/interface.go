package auth

import (
	"context"
	"geedbro-website-backend/internal/auth"
	"geedbro-website-backend/internal/entity"
	"geedbro-website-backend/internal/pkg"
	"geedbro-website-backend/internal/repository/postgres/user"
)

type Auth interface {
	GenerateToken(ctx context.Context, data auth.GenerateToken) (string, error)
	IsValidToken(ctx context.Context, token string) (entity.User, error)
	GetTokenData(ctx context.Context, token string) (auth.TokenData, error)
}

type User interface {
	GetByUsername(ctx context.Context, username string) (user.AdminGetDetail, *pkg.Error)
}
