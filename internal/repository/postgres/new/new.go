package new

import (
	"context"
	"fmt"
	"internal/entity"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"gitlab.com/blogs/internal/pkg"
	"gitlab.com/blogs/internal/pkg/repository/postgres"
)

type Repository struct {
	*postgres.Database
}

func NewRepository(postgresDB *postgres.Database) *Repository {
	return &Repository{postgresDB}
}

func (r Repository) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetListResponse, int, *pkg.Error) {
	query := fmt.Sprintf(`
		SELECT
			id,
			title,
			content,
			source
		FROM
		    news
		WHERE deleted_at IS NULL
	`)

	whereNews := ""

	if filter.Title != nil {
		title := strings.Replace(*filter.Title, " ", "", -1)
		whereNews += fmt.Sprintf(" AND REPLACE(title, ' ', '') ilike '%s'", "%"+title+"%")
	}
	query += whereNews

	if filter.Limit != nil {
		query += fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		query += fmt.Sprintf("OFFSET %d", *filter.Offset)
	}
	rows, er := r.QueryContext(ctx, query)
	if er != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(er, "selecting user list"),
			Status: http.StatusInternalServerError,
		}
	}
	var list []AdminGetListResponse
	for rows.Next() {
		var detail AdminGetListResponse
		if er = rows.Scan(&detail.Id, &detail.Title, &detail.Content, &detail.Source); er != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(er, "scanning new"),
				Status: http.StatusInternalServerError,
			}
		}
		list = append(list, detail)
	}
	countQuery := fmt.Sprintf(`
	SELECT
	COUNT(*)
	FROM
		news
	WHERE deleted_at IS NULL
`)
	countRows, er := r.QueryContext(ctx, countQuery+whereNews)
	if er != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(er, "selecting news count"),
			Status: http.StatusInternalServerError,
		}
	}
	count := 0

	for countRows.Next() {
		if er = countRows.Scan(&count); er != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(er, "scanning news count"),
				Status: http.StatusInternalServerError,
			}
		}
	}
	return list, count, nil
}

func (r Repository) AdminGetById(ctx context.Context, id string) (AdminGetDetail, *pkg.Error) {
	var detail AdminGetDetail

	err := r.NewSelect().Model(&detail).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return AdminGetDetail{}, &pkg.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, *pkg.Error) {
	var response AdminCreateResponse

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return AdminCreateResponse{}, er
	}
	if err := r.ValidateStruct(&request, "Title", "Content", "Source"); err != nil {
		return AdminCreateResponse{}, err
	}

	response.Id = uuid.NewString()
	response.Title = request.Title
	response.Content = request.Content
	response.Source = request.Source
	response.CreatedBy = &dataCtx.UserId
	response.CreatedAt = time.Now()
	err := r.ManualInsert(ctx, &response, "AdminCreate")
	if err != nil {
		return AdminCreateResponse{}, err
	}

	return response, nil
}

func (r Repository) AdminUpdate(ctx context.Context, request AdminUpdateRequest) *pkg.Error {
	userData, err := r.AdminGetById(ctx, request.Id)
	if err != nil {
		return err
	}
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	q := r.NewUpdate().Table("news").Where("deleted_at is null AND id = ?", request.Id)

	if request.Title != nil {
		q.Set("title = ?", request.Title)

	}
	if request.Content != nil {
		q.Set("content = ?", request.Content)

	}
	if request.Source != nil {
		q.Set("source = ?", request.Source)

	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", dataCtx.UserId)

	_, err1 := q.Exec(ctx)
	if err1 != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err1, "updating user"),
			Status: http.StatusInternalServerError,
		}
	}
	newUpdateData, err := r.AdminGetById(ctx, request.Id)
	if err != nil {
		return err
	}
	updateData := map[string]interface{}{
		"oldData": userData,
		"newData": newUpdateData,
	}

	var loggerData entity.LogCreateDto
	loggerData.Action = "AdminUpdateAll"
	loggerData.Method = "PUT"
	loggerData.Data = updateData
	err2 := r.LogCreate(ctx, loggerData)
	if err2 != nil {
		return err2
	}
	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id, role string) *pkg.Error {

	return r.DeleteRow(ctx, "news", id, role)
}
