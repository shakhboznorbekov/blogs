package new

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/blogs/internal/pkg"
	new2 "github.com/blogs/internal/repository/postgres/new"
	"github.com/blogs/internal/service/request"
	"github.com/blogs/internal/service/response"
)

type Controller struct {
	new New
}

func NewController(new New) *Controller {
	return &Controller{new}
}

// AdminGetNewList godoc
// @Security ApiKeyAuth
// @Summary Get New List
// @Description  Get New List
// @Tags New
// @Accept json
// @Produce json
// @Param offset query integer false "offset"
// @Param limit query integer false "limit"
// @Param title query string false "title"
// @Success 200 {object} new.AdminGetListResponse
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/news/list [GET]
func (cl Controller) AdminGetNewList(c *gin.Context) {
	var filter new2.Filter
	fieldErrors := make([]pkg.FieldError, 0)

	limit, err := request.GetQuery(c, reflect.Int, "limit")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := limit.(*int); ok {
		filter.Limit = value
	}

	offset, err := request.GetQuery(c, reflect.Int, "offset")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := offset.(*int); ok {
		filter.Offset = value
	}

	title, err := request.GetQuery(c, reflect.String, "title")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := title.(*string); ok {
		filter.Title = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	data, count, er := cl.new.AdminGetList(c, filter)
	if er != nil {
		response.RespondError(c, er)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data": map[string]interface{}{
			"results": data,
			"count":   count,
		},
	})
}

// AdminGetNewDetail godoc
// @Security ApiKeyAuth
// @Summary Get New ById
// @Description  Get New ById
// @Tags New
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} new.AdminDetailResponseSwagger
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/news/{id} [GET]
func (cl Controller) AdminGetNewDetail(c *gin.Context) {
	idParam, err := request.GetParam(c, reflect.String, "id")
	var id string
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else if value, ok := idParam.(string); ok {
		id = value
	}

	data, er := cl.new.AdminGetById(c, id)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	response.Respond(c, gin.H{
		"status": true,
		"data":   data,
	})
}

// AdminCreateNew godoc
// @Security ApiKeyAuth
// @Summary  New
// @Description  Create New
// @Tags New
// @Accept json
// @Produce json
// @Param new body new.AdminCreateRequest true "new"
// @Success 200 {object} new.AdminCreateResponseSwagger
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/news/create [POST]
func (cl Controller) AdminCreateNew(c *gin.Context) {
	var data new2.AdminCreateRequest

	er := request.BindFunc(c, &data)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	detail, er := cl.new.AdminCreate(c, data)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

// AdminUpdateNew godoc
// @Security ApiKeyAuth
// @Summary Update New
// @Description  Update New
// @Tags New
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param new body  new.AdminUpdateRequest true "new"
// @Success 200 {object} response.StatusOk
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/news/{id} [PUT]
func (cl Controller) AdminUpdateNew(c *gin.Context) {
	idParam, err := request.GetParam(c, reflect.String, "id")
	var id string
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else if value, ok := idParam.(string); ok {
		id = value
	}
	var data new2.AdminUpdateRequest

	er := request.BindFunc(c, &data)
	if er != nil {
		c.JSON(er.Status, gin.H{
			"message": er.Err.Error(),
			"status":  false,
		})

		return
	}
	if data.Id == "" {
		data.Id = id
	}

	err2 := cl.new.AdminUpdate(c, data)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err2.Err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}

// AdminDeleteNew godoc
// @Security ApiKeyAuth
// @Summary  Delete New
// @Description  Delete New
// @Tags New
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} response.StatusOk
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/news/{id} [DELETE]
func (cl Controller) AdminDeleteNew(c *gin.Context) {
	idParam, err1 := request.GetParam(c, reflect.String, "id")
	var id string
	if err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
	} else if value, ok := idParam.(string); ok {
		id = value
	}

	err := cl.new.AdminDelete(c, id, "Admin")
	if err != nil {
		response.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}
