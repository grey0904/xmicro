package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xmicro/internal/constant"
	"xmicro/internal/httperror"
	"xmicro/internal/utils/u_jwt"
)

type UserController struct{}

func (that *UserController) Login(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		err := httperror.NewHttpError("Query parameter not found", "name query parameter is required", http.StatusBadRequest)
		c.Error(err)
		return
	}

	token, err := u_jwt.CreateToken(
		u_jwt.Claims{
			Id:   "1",
			Name: username,
		},
		[]byte(constant.JwtSalt))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, token)
}
