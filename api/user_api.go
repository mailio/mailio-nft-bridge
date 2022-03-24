package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/service"
)

type UserAPI struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUserAPI(service *service.UserService) *UserAPI {
	return &UserAPI{
		service:  service,
		validate: validator.New(),
	}
}

// Admin login godoc
// @Summary      Admin login
// @Description  User login and returns JWT token
// @Description  User considered as admin, since only adding catalogs is allowed
// @Tags         Auth
// @Param        user  body      model.EmailPasswordInput  true  "Email and Password required"
// @Success      200   {object}  model.JwtTokenOutput
// @Failure      401   {object}  api.JSONError  "login failed"
// @Failure      500   {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/login [post]
func (u *UserAPI) Login(c *gin.Context) {
	var input model.EmailPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := u.validate.Struct(input); err != nil {
		AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := u.service.Login(input.Email, input.Password)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
