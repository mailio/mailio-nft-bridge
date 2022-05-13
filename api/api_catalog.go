package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/service"
)

type NftCatalogAPI struct {
	service  *service.NftCatalogService
	validate *validator.Validate
}

func NewNftCatalogAPI(service *service.NftCatalogService) *NftCatalogAPI {
	return &NftCatalogAPI{
		service:  service,
		validate: validator.New(),
	}
}

// Get Catalog
// @Summary      Get Catalog
// @Description  Get Catalog by id
// @Tags         Catalog
// @Param        id   path      string  true  "id"
// @Success      200  {object}  model.Catalog
// @Failure      404  {object}  api.JSONError  "catalog not found"
// @Failure      500    {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/catalog/{id} [get]
func (ca *NftCatalogAPI) GetCatalog(c *gin.Context) {
	id := c.Param("id")
	cat, err := ca.service.GetCatalog(id)
	if err != nil {
		if err == model.ErrNotFound {
			AbortWithError(c, http.StatusNotFound, "catalog not found")
			return
		}
		AbortWithError(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, cat)
}

// Upsert catalog
// @Security     ApiKeyAuth
// @Summary      Upsert Catalog
// @Description  When ID is given with the POST object then it's an update, otherwise insert
// @Tags         Catalog
// @Param        catalog  body      model.Catalog  true  "catalog"
// @Success      200      {object}  model.Catalog
// @Failure      400      {object}  api.JSONError  "invalid input"
// @Failure      401      {object}  api.JSONError  "access denied"
// @Failure      500      {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/catalog [post]
func (ca *NftCatalogAPI) PutCatalog(c *gin.Context) {
	cat := &model.Catalog{}
	if err := c.ShouldBindJSON(cat); err != nil {
		AbortWithError(c, http.StatusBadRequest, "invalid json body")
		return
	}
	err := ca.validate.Struct(cat)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	cat, err = ca.service.PutCatalog(cat)
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, cat)
}

// List Catalogs
// @Summary      List Catalog
// @Description  List Catalogs
// @Tags         Catalog
// @Param        limit  query     int  false  "limit"
// @Success      200    {array}   model.Catalog
// @Failure      500  {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/catalog [get]
func (ca *NftCatalogAPI) ListCatalogs(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 500
	if limitStr != "" {
		l, cErr := strconv.Atoi(limitStr)
		if cErr != nil {
			AbortWithError(c, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = l
	}
	cats, err := ca.service.ListAllCatalogs(limit)
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, cats)
}
