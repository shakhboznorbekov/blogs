package blog

import (
	"context"
	"fmt"
	"github.com/blogs/internal/entity"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/blogs/internal/pkg"
	"github.com/blogs/internal/pkg/repository/postgres"
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
			author
		FROM
		    blogs
		WHERE deleted_at IS NULL
	`)

	whereTitle := ""

	if filter.Title != nil {
		title := strings.Replace(*filter.Title, " ", "", -1)
		whereTitle += fmt.Sprintf(" AND REPLACE(title, ' ', '') ilike '%s'", "%"+title+"%")
	}
	query += whereTitle

	if filter.Limit != nil {
		query += fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		query += fmt.Sprintf("OFFSET %d", *filter.Offset)
	}
	rows, er := r.QueryContext(ctx, query)
	if er != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(er, "selecting blog list"),
			Status: http.StatusInternalServerError,
		}
	}
	var list []AdminGetListResponse
	for rows.Next() {
		var detail AdminGetListResponse
		if er = rows.Scan(&detail.Id, &detail.Title, &detail.Content, &detail.Author); er != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(er, "scanning blog"),
				Status: http.StatusInternalServerError,
			}
		}
		list = append(list, detail)
	}
	countQuery := fmt.Sprintf(`
	SELECT
	COUNT(*)
	FROM
		blogs
	WHERE deleted_at IS NULL
`)
	countRows, er := r.QueryContext(ctx, countQuery+whereTitle)
	if er != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(er, "selecting blogs count"),
			Status: http.StatusInternalServerError,
		}
	}
	count := 0

	for countRows.Next() {
		if er = countRows.Scan(&count); er != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(er, "scanning blogs count"),
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
	if err := r.ValidateStruct(&request, "Title", "Content", "Author"); err != nil {
		return AdminCreateResponse{}, err
	}

	response.Id = uuid.NewString()
	response.Title = request.Title
	response.Content = request.Content
	response.Author = request.Author
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
	q := r.NewUpdate().Table("blogs").Where("deleted_at is null AND id = ?", request.Id)

	if request.Title != nil {
		q.Set("title = ?", request.Title)

	}
	if request.Content != nil {
		q.Set("content = ?", request.Content)

	}
	if request.Author != nil {
		q.Set("author = ?", request.Author)

	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", dataCtx.UserId)

	_, err1 := q.Exec(ctx)
	if err1 != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err1, "updating blog"),
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

	return r.DeleteRow(ctx, "blogs", id, role)
}
