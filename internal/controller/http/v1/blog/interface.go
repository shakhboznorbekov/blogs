package blog

import (
	"context"
	"gitlab.com/blogs/internal/repository/postgres/blog"

	"gitlab.com/blogs/internal/pkg"
)

type Blog interface {
	AdminGetList(ctx context.Context, filter blog.Filter) ([]blog.AdminGetListResponse, int, *pkg.Error)
	AdminGetById(ctx context.Context, id string) (blog.AdminGetDetail, *pkg.Error)
	AdminCreate(ctx context.Context, request blog.AdminCreateRequest) (blog.AdminCreateResponse, *pkg.Error)
	AdminUpdate(ctx context.Context, request blog.AdminUpdateRequest) *pkg.Error
	AdminDelete(ctx context.Context, id, role string) *pkg.Error
}
