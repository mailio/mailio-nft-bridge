package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/service"
)

type NftImagesAPI struct {
	service *service.NftImagesService
}

func NewNftImagesAPI(service *service.NftImagesService) *NftImagesAPI {
	return &NftImagesAPI{
		service: service,
	}
}

// Nft images list
// @Tags         Nft Images
// @Security     ApiKeyAuth
// @Summary      List all pinned images
// @Description  List all images pinned to infura
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.NftImage
// @Failure      400  {object}  api.JSONError  "failed to read file"
// @Failure      500  {object}  api.JSONError  "upload failed"
// @Router       /v1/nftimage/list [get]
func (ni *NftImagesAPI) List(c *gin.Context) {
	list, err := ni.service.List()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("failed to list images"))
		return
	}
	c.JSON(http.StatusOK, list)
}

// Nft images upload
// @Tags         Nft Images
// @Security     ApiKeyAuth
// @Summary      Upload to IPFS
// @Description  Upload file to IPFS (on Infura) and pin it
// @ID           file.upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "image file"
// @Success      200    {object}  model.NftImageUploadResponse
// @Failure      400    {object}  api.JSONError  "failed to read file"
// @Failure      500    {object}  api.JSONError  "upload failed"
// @Router       /v1/nftimage/upload [post]
func (ni *NftImagesAPI) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		lc.Log.Error("failed to store image. image required", err)
		c.AbortWithError(http.StatusBadRequest, errors.New("no image received"))
		return
	}
	resp, err := ni.service.Upload(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("failed to upload image"))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Nft delete image
// @Tags         Nft Images
// @Security     ApiKeyAuth
// @Summary      Delete pinned image
// @Description  Delete pinned image from local service (infura)
// @Param        hash  path  string  true  "image hash"
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.NftPins
// @Failure      400  {object}  api.JSONError  "failed to read file"
// @Failure      500  {object}  api.JSONError  "upload failed"
// @Router       /v1/nftimage/{hash} [delete]
func (ni *NftImagesAPI) RemovePin(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("no hash received"))
		return
	}
	removedPins, err := ni.service.RemovePin(hash)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("failed to remove pin"))
		return
	}
	c.JSON(http.StatusOK, removedPins)
}
