package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"gitlab.com/admin_service_shifo/internal/entity"
	"gitlab.com/admin_service_shifo/internal/pkg"
	"gitlab.com/admin_service_shifo/internal/pkg/repository/postgres"
	"gitlab.com/admin_service_shifo/internal/service/hash"
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
			username,
			status,
			role
		FROM
		    users
		WHERE deleted_at IS NULL
	`)

	whereUser := ""

	if filter.Username != nil {
		username := strings.Replace(*filter.Username, " ", "", -1)
		whereUser += fmt.Sprintf(" AND REPLACE(username, ' ', '') ilike '%s'", "%"+username+"%")
	}
	query += whereUser

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
		if er = rows.Scan(&detail.Id, &detail.Username, &detail.Status, &detail.Role); er != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(er, "scanning user"),
				Status: http.StatusInternalServerError,
			}
		}
		list = append(list, detail)
	}
	countQuery := fmt.Sprintf(`
	SELECT
	COUNT(*)
	FROM
		users
	WHERE deleted_at IS NULL
`)
	countRows, er := r.QueryContext(ctx, countQuery+whereUser)
	if er != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(er, "selecting user count"),
			Status: http.StatusInternalServerError,
		}
	}
	count := 0

	for countRows.Next() {
		if er = countRows.Scan(&count); er != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(er, "scanning user count"),
				Status: http.StatusInternalServerError,
			}
		}
	}
	fmt.Println(*list[0].Status)
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

func (r Repository) GetByUsername(ctx context.Context, username string) (AdminGetDetail, *pkg.Error) {
	var detail AdminGetDetail
	err := r.NewSelect().Model(&detail).Where("username = ?", username).Scan(ctx)

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
	if err := r.ValidateStruct(&request, "Username", "Password"); err != nil {
		return AdminCreateResponse{}, err
	}

	response.Id = uuid.NewString()
	response.Username = request.Username
	role := "User"
	response.Role = &role
	hashPassword, err2 := hash.HashPassword(request.Password)
	if err2 != nil {
		return AdminCreateResponse{}, &pkg.Error{
			Err:    err2,
			Status: http.StatusInternalServerError,
		}
	}
	response.Status = request.Status
	response.Password = hashPassword
	response.CreatedBy = &dataCtx.UserId
	response.CreatedAt = time.Now()
	err := r.ManualInsert(ctx, &response, "AdminCreate")
	if err != nil {
		return AdminCreateResponse{}, err
	}

	return response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) *pkg.Error {
	userData, err := r.AdminGetById(ctx, request.Id)
	if err != nil {
		return err
	}
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	q := r.NewUpdate().Table("users").Where("deleted_at is null AND id = ?", request.Id)

	if request.Username != nil {
		q.Set("username = ?", request.Username)

	}
	if request.Status != nil {
		q.Set("status = ?", request.Status)

	}
	if request.Role != nil {
		q.Set("role = ?", request.Role)

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

func (r Repository) AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) *pkg.Error {
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	if err := r.ValidateStruct(&request, "Id"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("users").Where("deleted_at is null AND id = ?", request.Id)

	if request.Username != nil {
		q.Set("username = ?", request.Username)
	}
	if request.Role != nil {
		q.Set("role = ?", request.Role)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", dataCtx.UserId)

	_, er1 := q.Exec(ctx)
	if er1 != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(er1, "updating user"),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id, role string) *pkg.Error {

	return r.DeleteRow(ctx, "users", id, role)
}
