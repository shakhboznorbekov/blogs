package blog

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blogs/internal/pkg/repository/postgres"
	blog_repo "github.com/blogs/internal/repository/postgres/blog"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestAdminGetList(t *testing.T) {
	// Create a new mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := bun.NewDB(mockDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	postgresDB := &postgres.Database{DB: db}

	// Create your repository with the mock DB
	repo := blog_repo.NewRepository(postgresDB)

	// Define your test data
	expectedList := []AdminGetListResponse{
		{Id: "72679d9e-7c4c-4453-bc49-faf19f7ae427", Title: "Blog 1", Content: "Content 1", Author: "Author 1"},
		{Id: "fc9c7f8d-34e4-4d00-a2e4-640b33a73b4d", Title: "Blog 2", Content: "Content 2", Author: "Author 2"},
	}
	expectedCount := 2

	// Set expectations for the mock DB
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author"}).
		AddRow(expectedList[0].Id, expectedList[0].Title, expectedList[0].Content, expectedList[0].Author).
		AddRow(expectedList[1].Id, expectedList[1].Title, expectedList[1].Content, expectedList[1].Author)

	mock.ExpectQuery("SELECT id, title, content, author FROM blogs WHERE deleted_at IS NULL").
		WillReturnRows(rows)

	countRows := sqlmock.NewRows([]string{"count"}).
		AddRow(expectedCount)

	mock.ExpectQuery("SELECT COUNT(*) FROM blogs WHERE deleted_at IS NULL").
		WillReturnRows(countRows)

	// Call the method being tested
	filter := Filter{Title: sql.NullString{String: "Blog", Valid: true}}
	actualList, actualCount, _ := repo.AdminGetList(context.Background(), filter)

	// Assert the expected result
	assert.Equal(t, expectedList, actualList)
	assert.Equal(t, expectedCount, actualCount)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAdminGetById(t *testing.T) {
	// Create a new mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := bun.NewDB(mockDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	postgresDB := &postgres.Database{DB: db}

	// Create your repository with the mock DB
	repo := blog_repo.NewRepository(postgresDB)

	// Define your test data
	expectedDetail := AdminGetDetail{
		Id:      "fc9c7f8d-34e4-4d00-a2e4-640b33a73b4d",
		Title:   "Blog 1",
		Content: "Content 1",
		Author:  "Author 1",
	}

	// Set expectations for the mock DB
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author"}).
		AddRow(expectedDetail.Id, expectedDetail.Title, expectedDetail.Content, expectedDetail.Author)

	mock.ExpectQuery("SELECT id, title, content, author FROM yourtable WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	// Call the method being tested
	actualDetail, _ := repo.AdminGetById(context.Background(), "fc9c7f8d-34e4-4d00-a2e4-640b33a73b4d")

	// Assert the expected result
	assert.Equal(t, expectedDetail, actualDetail)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAdminCreate(t *testing.T) {
	// Create a new mock DB

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := bun.NewDB(mockDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	postgresDB := &postgres.Database{DB: db}

	// Create your repository with the mock DB
	repo := blog_repo.NewRepository(postgresDB)

	// Define your test data
	request := AdminCreateRequest{
		Title:   "New Blog",
		Content: "Blog Content",
		Author:  "Shaxboz",
	}

	expectedResponse := AdminCreateResponse{
		Id:        "generated_uuid",
		Title:     request.Title,
		Content:   request.Content,
		Author:    request.Author,
		CreatedBy: &request.Title,
		CreatedAt: time.Now(),
	}

	// Set expectations for the mock DB
	mock.ExpectQuery("INSERT INTO blogs").
		WithArgs(request.Title, request.Content, request.Author, expectedResponse.CreatedBy, expectedResponse.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedResponse.Id))

	// Call the method being tested
	actualResponse, _ := repo.AdminCreate(context.Background(), request)

	// Assert the expected result
	assert.Equal(t, expectedResponse, actualResponse)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAdminUpdate(t *testing.T) {
	// Create a new mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := bun.NewDB(mockDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Create your repository with the mock DB
	repo := &postgres.Database{DB: db}

	// Define your test data
	request := AdminUpdateRequest{
		Id:      "fc9c7f8d-34e4-4d00-a2e4-640b33a73b4d",
		Title:   "Updated Title",
		Content: "Updated Content",
		Author:  "Updated Author",
	}

	expectedOldData := AdminGetDetail{
		Id:      "fc9c7f8d-34e4-4d00-a2e4-640b33a73b4d",
		Title:   "Original Title",
		Content: "Original Content",
		Author:  "Original Author",
	}

	expectedNewData := AdminGetDetail{
		Id:      "fc9c7f8d-34e4-4d00-a2e4-640b33a73b4d",
		Title:   "Updated Title",
		Content: "Updated Content",
		Author:  "Updated Author",
	}

	// Set expectations for the mock DB
	mock.ExpectQuery("SELECT id, title, content, author FROM blogs WHERE id = ?").
		WithArgs(request.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author"}).
			AddRow(expectedOldData.Id, expectedOldData.Title, expectedOldData.Content, expectedOldData.Author))

	mock.ExpectExec("UPDATE yourtable SET title = ?, content = ?, author = ?, updated_at = ?, updated_by = ? WHERE deleted_at is null AND id = ?").
		WithArgs(request.Title, request.Content, request.Author, sqlmock.AnyArg(), sqlmock.AnyArg(), request.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("SELECT id, title, content, author FROM blogs WHERE id = ?").
		WithArgs(request.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author"}).
			AddRow(expectedNewData.Id, expectedNewData.Title, expectedNewData.Content, expectedNewData.Author))

	mock.ExpectQuery("INSERT INTO blogs").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("generated_uuid"))

	// Call the method being tested
	err = repo.AdminUpdate(context.Background(), request)

	// Assert the expected result
	assert.NoError(t, err)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAdminDelete(t *testing.T) {
	// Open a real PostgreSQL connection for the test

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := bun.NewDB(mockDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Create your repository with the mock DB
	repo := &postgres.Database{DB: db}

	// Set expectations for the mock DB
	mock.ExpectExec("UPDATE blogs SET deleted_at = NOW(), deleted_by = ? WHERE deleted_at is null AND id = ?").
		WithArgs("admin", "1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method being tested with the mock DB
	err = repo.AdminDelete(context.Background(), "1", "admin")

	// Assert the expected result
	assert.NoError(t, err)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

// ConnectDB - It will be used anywhere in application to connect database.
func ConnectDB() (*gorm.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("DATABASE_NAME")

	db, err := gorm.Open("postgres", fmt.Sprintf("%v/%v", databaseURL, databaseName))
	if err != nil {
		log.Printf("Error connecting database.\n%v", err)
		return nil, err
	}

	return db, nil
}

// MockedDB is used in unit tests to mock db
func MockedDB(operation string) {
	/*
	   If tests are running in CI, environment variables should not be loaded.
	   The reason is environment vars will be provided through CI config file.
	*/
	if CI := os.Getenv("CI"); CI == "" {
		// If tests are not running in CI, we have to load .env file.
		_, fileName, _, _ := runtime.Caller(0)
		currPath := filepath.Dir(fileName)
		// path should be relative path from this directory to ".env"
		err := godotenv.Load(currPath + "/../../.env")
		if err != nil {
			log.Fatalf("Error loading env.\n%v", err)
		}
	}

	dbName := os.Getenv("DATABASE_NAME")
	pgUser := os.Getenv("PSQL_USER")
	pgPassword := os.Getenv("PSQL_PASSWORD")

	// createdb => https://www.postgresql.org/docs/7.0/app-createdb.htm
	// dropdb => https://www.postgresql.org/docs/7.0/app-dropdb.htm
	var command string

	if operation == CREATE {
		command = "createdb"
	} else {
		command = "dropdb"
	}

	// createdb & dropdb commands have same configuration syntax.
	cmd := exec.Command(command, "-h", "localhost", "-U", pgUser, "-e", dbName)
	cmd.Env = os.Environ()

	/*
	   if we normally execute createdb/dropdb, we will be propmted to provide password.
	   To inject password automatically, we have to set PGPASSWORD as prefix.
	*/
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", pgPassword))

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error executing %v on %v.\n%v", command, dbName, err)
	}

	/*
	   Alternatively instead of createdb/dropdb, you can use
	   psql -c "CREATE/DROP DATABASE DBNAME" "DATABASE_URL"
	*/
}
