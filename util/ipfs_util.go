package util

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/go-resty/resty/v2"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
)

// upload and PIN the image to IPFS
// TODO: add another pin to pinata cloud
func UploadToIPFSAndPin(ipfsClient *resty.Client, filename string, data io.Reader) ([]*model.NftImageUploadResponse, error) {

	response, respErr := ipfsClient.R().
		SetFileReader("file", filename, data).
		SetContentLength(true).
		Post("/api/v0/add")

	if respErr != nil {
		lc.Log.Error("Failed to upload file", respErr)
		return nil, respErr
	}
	respBody := response.Body()
	if response.StatusCode() != 200 {
		lc.Log.Error("Failed to upload file", string(respBody))
		return nil, errors.New("Failed to upload file to ipfs")
	}
	output := []*model.NftImageUploadResponse{}
	splitted := strings.Split(string(respBody), "\n")
	for _, split := range splitted {
		if len(split) > 3 {
			var uploadResponse model.NftImageUploadResponse
			mErr := json.Unmarshal([]byte(split), &uploadResponse)
			if mErr != nil {
				lc.Log.Error("failed to unmarshal upload response", string(respBody), mErr)
				return nil, mErr
			}
			output = append(output, &uploadResponse)
		}
	}
	return output, nil
}
