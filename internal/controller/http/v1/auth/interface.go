package auth

import (
	"context"
	"github.com/blogs/internal/auth"
	"github.com/blogs/internal/entity"
	"github.com/blogs/internal/pkg"
	"github.com/blogs/internal/repository/postgres/user"
)

type Auth interface {
	GenerateToken(ctx context.Context, data auth.GenerateToken) (string, error)
	IsValidToken(ctx context.Context, token string) (entity.User, error)
	GetTokenData(ctx context.Context, token string) (auth.TokenData, error)
}

type User interface {
	GetByUsername(ctx context.Context, username string) (user.AdminGetDetail, *pkg.Error)
}
