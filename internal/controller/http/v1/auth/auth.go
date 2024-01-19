package auth

import (
	auth_service "github.com/blogs/internal/auth"
	"github.com/blogs/internal/service/hash"
	"github.com/blogs/internal/service/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	user User
	auth Auth
}

func NewController(user User, auth Auth) *Controller {
	return &Controller{user, auth}
}

// SignIn godoc
// @Summary Sign In
// @Description  Create Author
// @Tags auth
// @Accept json
// @Produce json
// @Param profile body auth.SignIn true "CreateAuthorRequestBody"
// @Success 200 {object} auth.AuthResponse
// @Response 400 {object} string "Invalid argument"
// @Failure 500 {object} string "Server Error"
// @Router /api/v1/user/sign-in [POST]
func (ac Controller) SignIn(c *gin.Context) {
	var data auth_service.SignIn

	if er := request.BindFunc(c, &data, "Username", "Password"); er != nil {
		c.JSON(er.Status, gin.H{
			"message": er.Err.Error(),
			"status":  false,
		})

		return
	}

	userDetail, err := ac.user.GetByUsername(c, data.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "user not found",
			"status":  false,
		})
		return
	}
	answer := hash.CheckPasswordHash(data.Password, *userDetail.Password)
	if answer == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "incorrect password!",
			"status":  false,
		})
		return
	}

	var generateTokenData auth_service.GenerateToken

	generateTokenData.Username = userDetail.Username
	generateTokenData.Role = userDetail.Role
	token, err2 := ac.auth.GenerateToken(c, generateTokenData)

	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err2.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
