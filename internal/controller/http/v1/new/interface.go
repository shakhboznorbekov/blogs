package new

import (
	"context"
	"github.com/blogs/internal/pkg"
	"github.com/blogs/internal/repository/postgres/new"
)

type New interface {
	AdminGetList(ctx context.Context, filter new.Filter) ([]new.AdminGetListResponse, int, *pkg.Error)
	AdminGetById(ctx context.Context, id string) (new.AdminGetDetail, *pkg.Error)
	AdminCreate(ctx context.Context, request new.AdminCreateRequest) (new.AdminCreateResponse, *pkg.Error)
	AdminUpdate(ctx context.Context, request new.AdminUpdateRequest) *pkg.Error
	AdminDelete(ctx context.Context, id, role string) *pkg.Error
}
