package response

import (
	"gitlab.com/blogs/internal/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func RespondError(c *gin.Context, err *pkg.Error) {
	c.JSON(err.Status, map[string]interface{}{
		"error":  err.Err.Error(),
		"f":      err.Fields,
		"status": false,
	})
}

type StatusOk struct {
	Massage bool `json:"massage"`
}
