package user

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/blogs/internal/pkg"
	user2 "github.com/blogs/internal/repository/postgres/user"
	"github.com/blogs/internal/service/request"
	"github.com/blogs/internal/service/response"
)

type Controller struct {
	user User
	auth Auth
}

func NewController(user User, auth Auth) *Controller {
	return &Controller{user, auth}
}

// AdminGetUserList godoc
// @Security ApiKeyAuth
// @Summary Get User List
// @Description  Get User List
// @Tags User
// @Accept json
// @Produce json
// @Param offset query integer false "offset"
// @Param limit query integer false "limit"
// @Param username query string false "username"
// @Success 200 {object} user.AdminGetListResponse
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/user/list [GET]
func (cl Controller) AdminGetUserList(c *gin.Context) {
	var filter user2.Filter
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

	username, err := request.GetQuery(c, reflect.String, "username")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := username.(*string); ok {
		filter.Username = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	data, count, er := cl.user.AdminGetList(c, filter)
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

// AdminGetUserDetail godoc
// @Security ApiKeyAuth
// @Summary Get User ById
// @Description  Get User ById
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} user.AdminDetailResponseSwagger
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/user/{id} [GET]
func (cl Controller) AdminGetUserDetail(c *gin.Context) {
	idParam, err := request.GetParam(c, reflect.String, "id")
	var id string
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else if value, ok := idParam.(string); ok {
		id = value
	}

	data, er := cl.user.AdminGetById(c, id)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	response.Respond(c, gin.H{
		"status": true,
		"data":   data,
	})
}

// AdminCreateUser godoc
// @Security ApiKeyAuth
// @Summary  User
// @Description  Create User
// @Tags User
// @Accept json
// @Produce json
// @Param user body user.AdminCreateRequest true "user"
// @Success 200 {object} user.AdminCreateResponseSwagger
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/user/create [POST]
func (cl Controller) AdminCreateUser(c *gin.Context) {
	var data user2.AdminCreateRequest

	er := request.BindFunc(c, &data)
	if er != nil {
		response.RespondError(c, er)

		return
	}

	detail, er := cl.user.AdminCreate(c, data)
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

// AdminUpdateUser godoc
// @Security ApiKeyAuth
// @Summary Update User
// @Description  Update User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param user body  user.AdminUpdateRequest true "user"
// @Success 200 {object} response.StatusOk
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/user/{id} [PUT]
func (cl Controller) AdminUpdateUser(c *gin.Context) {
	idParam, err := request.GetParam(c, reflect.String, "id")
	var id string
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else if value, ok := idParam.(string); ok {
		id = value
	}
	var data user2.AdminUpdateRequest

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

	err2 := cl.user.AdminUpdateAll(c, data)
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

// AdminDeleteUser godoc
// @Security ApiKeyAuth
// @Summary  Delete User
// @Description  Delete User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} response.StatusOk
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/admin/user/{id} [DELETE]
func (cl Controller) AdminDeleteUser(c *gin.Context) {
	idParam, err1 := request.GetParam(c, reflect.String, "id")
	var id string
	if err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
	} else if value, ok := idParam.(string); ok {
		id = value
	}

	err := cl.user.AdminDelete(c, id, "Admin")
	if err != nil {
		response.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}
