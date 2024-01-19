package blog

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/blogs/internal/pkg"
	blog2 "github.com/blogs/internal/repository/postgres/blog"
	"github.com/blogs/internal/service/request"
	"github.com/blogs/internal/service/response"
)

type Controller struct {
	blog Blog
}

func NewController(blog Blog) *Controller {
	return &Controller{blog}
}

// AdminGetBlogList godoc
// @Security ApiKeyAuth
// @Summary Get Blog List
// @Description  Get Blog List
// @Tags Blog
// @Accept json
// @Produce json
// @Param offset query integer false "offset"
// @Param limit query integer false "limit"
// @Param title query string false "title"
// @Success 200 {object} blog.AdminGetListResponse
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/blog/list [GET]
func (cl Controller) AdminGetBlogList(c *gin.Context) {
	var filter blog2.Filter
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

	data, count, er := cl.blog.AdminGetList(c, filter)
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

// AdminGetBlogDetail godoc
// @Security ApiKeyAuth
// @Summary Get Blog ById
// @Description  Get Blog ById
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} blog.AdminDetailResponseSwagger
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/blog/{id} [GET]
func (cl Controller) AdminGetBlogDetail(c *gin.Context) {
	idParam, err := request.GetParam(c, reflect.String, "id")
	var id string
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else if value, ok := idParam.(string); ok {
		id = value
	}

	data, er := cl.blog.AdminGetById(c, id)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	response.Respond(c, gin.H{
		"status": true,
		"data":   data,
	})
}

// AdminCreateBlog godoc
// @Security ApiKeyAuth
// @Summary  Blog
// @Description  Create Blog
// @Tags Blog
// @Accept json
// @Produce json
// @Param blog body blog.AdminCreateRequest true "blog"
// @Success 200 {object} blog.AdminCreateResponseSwagger
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/blog/create [POST]
func (cl Controller) AdminCreateBlog(c *gin.Context) {
	var data blog2.AdminCreateRequest

	er := request.BindFunc(c, &data)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	detail, er := cl.blog.AdminCreate(c, data)
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

// AdminUpdateBlog godoc
// @Security ApiKeyAuth
// @Summary Update Blog
// @Description  Update Blog
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param blog body  blog.AdminUpdateRequest true "blog"
// @Success 200 {object} response.StatusOk
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/blog/{id} [PUT]
func (cl Controller) AdminUpdateBlog(c *gin.Context) {
	idParam, err := request.GetParam(c, reflect.String, "id")
	var id string
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else if value, ok := idParam.(string); ok {
		id = value
	}
	var data blog2.AdminUpdateRequest

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

	err2 := cl.blog.AdminUpdate(c, data)
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

// AdminDeleteBlog godoc
// @Security ApiKeyAuth
// @Summary  Delete Blog
// @Description  Delete Blog
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} response.StatusOk
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/blog/{id} [DELETE]
func (cl Controller) AdminDeleteBlog(c *gin.Context) {
	idParam, err1 := request.GetParam(c, reflect.String, "id")
	var id string
	if err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
	} else if value, ok := idParam.(string); ok {
		id = value
	}

	err := cl.blog.AdminDelete(c, id, "Admin")
	if err != nil {
		response.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}
