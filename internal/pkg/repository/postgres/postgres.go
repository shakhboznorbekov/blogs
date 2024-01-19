package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blogs/internal/entity"
	"github.com/blogs/internal/pkg"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"net/http"
	"reflect"
	"time"
)

type CtxData struct {
	UserId string
	Role   string
}

type Database struct {
	*bun.DB
	DefaultLang   string
	ServerBaseUrl string
}

func New(DBUsername, DBPassword, DBPort, DBName, defaultLang, serverBaseUrl string) *Database {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", DBUsername, DBPassword, DBPort, DBName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return &Database{
		DB:            db,
		DefaultLang:   defaultLang,
		ServerBaseUrl: serverBaseUrl,
	}
}

func (d Database) DeleteRow(ctx context.Context, table, id, role string) *pkg.Error {
	userID, ok := ctx.Value("userId").(string)
	if !ok || (userID == "") {
		return &pkg.Error{
			Err:    errors.New("userId in context is required"),
			Status: http.StatusInternalServerError,
		}
	}

	ctxRole, ok := ctx.Value("role").(string)

	if !ok || (ctxRole == "") {
		return &pkg.Error{
			Err:    errors.New("role in context is required"),
			Status: http.StatusInternalServerError,
		}
	}

	if ctxRole != role {
		return &pkg.Error{
			Err:    errors.New(fmt.Sprintf("you have not permission to delete from table: %s", table)),
			Status: http.StatusInternalServerError,
		}
	}
	fmt.Println("DELETE DOC TEST", id, userID, table)
	_, err := d.NewUpdate().
		Table(table).
		Where("deleted_at is null AND id = ?", id).
		Set("deleted_at = ?", time.Now()).
		Set("deleted_by = ?", userID).
		Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return &pkg.Error{
			Err:    errors.New("delete row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}
	var loggerData entity.Logger
	loggerData.Action = "Delete"
	loggerData.Method = "Delete"
	loggerData.CreatedAt = time.Now()
	_, err = d.NewInsert().Model(&loggerData).Exec(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (d Database) ValidateStruct(s interface{}, requiredFields ...string) *pkg.Error {
	structVal := reflect.Value{}
	if reflect.Indirect(reflect.ValueOf(s)).Kind() == reflect.Struct {
		structVal = reflect.Indirect(reflect.ValueOf(s))
	} else {
		return &pkg.Error{
			Err:    errors.New("input object should be struct"),
			Status: http.StatusBadRequest,
		}
	}

	errFields := make([]pkg.FieldError, 0)

	structType := reflect.Indirect(reflect.ValueOf(s)).Type()
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()
		if !isSet {
			log.Print(isSet, fieldName, reflect.ValueOf(field))
			for _, f := range requiredFields {
				if f == fieldName {
					errFields = append(errFields, pkg.FieldError{
						Err:   errors.New("field is required!"),
						Field: fieldName,
					})
				}
			}
		}
	}

	if len(errFields) > 0 {
		return &pkg.Error{
			Err:    errors.New("required fields"),
			Fields: errFields,
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

func (d Database) CheckCtx(ctx context.Context) (CtxData, *pkg.Error) {
	fieldErrors := make([]pkg.FieldError, 0)
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		fieldErrors = append(fieldErrors, pkg.FieldError{
			Err:   errors.New("missing field in ctx"),
			Field: "user_id",
		})
	}
	ctxRole, ok := ctx.Value("role").(string)
	if !ok {
		fieldErrors = append(fieldErrors, pkg.FieldError{
			Err:   errors.New("missing field in ctx"),
			Field: "role",
		})
	}

	if len(fieldErrors) > 0 {
		return CtxData{}, &pkg.Error{
			Err:    errors.New("missing fields in context"),
			Fields: fieldErrors,
			Status: http.StatusInternalServerError,
		}
	}

	return CtxData{
		UserId: userId,
		Role:   ctxRole,
	}, nil
}

func (d Database) GetLang(ctx context.Context) string {
	return d.DefaultLang
}

func (d Database) ManualInsert(ctx context.Context, data interface{}, action string) *pkg.Error {
	_, err := d.NewInsert().Model(data).Exec(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	var loggerData entity.Logger
	loggerData.Action = action
	loggerData.Method = "POST"
	loggerData.Data = data
	loggerData.CreatedAt = time.Now()
	_, err = d.NewInsert().Model(&loggerData).Exec(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (d Database) LogCreate(ctx context.Context, request entity.LogCreateDto) *pkg.Error {
	var loggerData entity.Logger
	loggerData.Action = request.Action
	loggerData.Method = request.Method
	loggerData.Data = request.Data
	loggerData.CreatedAt = time.Now()
	_, err := d.NewInsert().Model(&loggerData).Exec(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}
