package service

import (
	"encoding/json"
	"errors"
	"mime/multipart"

	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/util"
)

type NftImagesService struct {
	environment *model.Environment
}

func NewNftImagesService(environment *model.Environment) *NftImagesService {

	return &NftImagesService{
		environment: environment,
	}
}

// list all pinned nft images
func (ni *NftImagesService) List() (*model.NftImage, error) {

	resp, err := ni.environment.IpfsInfuraClient.R().
		SetQueryParam("type", "all").
		Post("/api/v0/pin/ls")
	if err != nil {
		lc.Log.Error("failed to list pinned images", err)
		return nil, err
	}
	if resp.StatusCode() != 200 {
		lc.Log.Error("failed to list pinned images", string(resp.Body()))
		return nil, errors.New("failed to list pinned images")
	}
	var images model.NftImage
	mErr := json.Unmarshal(resp.Body(), &images)
	if mErr != nil {
		lc.Log.Error("failed to unmarshal images", mErr)
		return nil, mErr
	}
	return &images, nil
}

// Upload file to IPFS (on Infura) and pin it
func (ni *NftImagesService) Upload(file *multipart.FileHeader) ([]*model.NftImageUploadResponse, error) {
	mpFile, err := file.Open()
	defer mpFile.Close()

	if err != nil {
		lc.Log.Error("failed to read file", err)
		return nil, errors.New("failed to read file")
	}

	return util.UploadToIPFSAndPin(ni.environment.IpfsInfuraClient, file.Filename, mpFile)
}

// deletes the pinned object
func (ni *NftImagesService) RemovePin(hash string) (*model.NftPins, error) {
	resp, err := ni.environment.IpfsInfuraClient.R().
		SetQueryParam("arg", hash).
		Post("/api/v0/pin/rm")

	if err != nil {
		lc.Log.Error("failed to removed pinned image", err)
		return nil, err
	}
	if resp.StatusCode() != 200 {
		lc.Log.Error("failed to delete pinned images", string(resp.Body()))
		return nil, errors.New("failed to delete pinned image")
	}
	var pins model.NftPins
	mErr := json.Unmarshal(resp.Body(), &pins)
	if mErr != nil {
		lc.Log.Error("failed to unmarshal images", mErr)
		return nil, mErr
	}
	return &pins, nil
}
