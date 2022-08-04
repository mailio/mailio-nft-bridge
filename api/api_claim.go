package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/service"
)

type ClaimAPI struct {
	service             *service.NftClaimService
	catalogService      *service.NftCatalogService
	validate            *validator.Validate
	httpClientReCaptcha *resty.Client
}

func NewClaimAPI(service *service.NftClaimService, catalogService *service.NftCatalogService) *ClaimAPI {
	return &ClaimAPI{
		service:             service,
		catalogService:      catalogService,
		validate:            validator.New(),
		httpClientReCaptcha: resty.New().SetHostURL(lc.Conf.ReCaptchaV3.Host),
	}
}

// List Claims
// @Summary      List Claims
// @Security     ApiKeyAuth
// @Description  List Latest claims
// @Tags         Claiming
// @Param        limit  query     int  false  "limit"
// @Success      200    {array}   model.Claim
// @Failure      500    {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/claim [get]
func (ca *ClaimAPI) ListClaims(c *gin.Context) {

	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		l, cErr := strconv.Atoi(limitStr)
		if cErr != nil {
			AbortWithError(c, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = l
	}

	claims, err := ca.service.ListClaims(limit)
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	c.JSON(http.StatusOK, claims)
}

// SafeMint new NFT
// @Summary      Mint new NFT
// @Description  Mints the new NFT based on the category selected. All NFTs are on Polygon
// @Tags         Claiming
// @Param        claim  body      model.Claim  true  "eip-712 signed claim"
// @Success      200    {array}   model.Claim
// @Failure      403    {object}  api.JSONError  "captacha failed"
// @Failure      400    {object}  api.JSONError  "invalid input"
// @Failure      500    {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/claim [post]
func (ca *ClaimAPI) MintClaim(c *gin.Context) {
	claim := &model.Claim{}
	if err := c.ShouldBindJSON(claim); err != nil {
		AbortWithError(c, http.StatusBadRequest, "invalid json body")
		return
	}
	err := ca.validate.Struct(claim)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	catalog, cErr := ca.catalogService.GetCatalog(claim.CatalogId)
	if cErr != nil {
		AbortWithError(c, http.StatusBadRequest, "Catalog invalid")
		return
	}

	// validate captcha v3
	reCaptchaResp, captchaErr := ca.httpClientReCaptcha.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetHeader("Cache-Control", "no-cache").
		SetFormData(map[string]string{
			"secret":   lc.Conf.ReCaptchaV3.Secret,
			"response": claim.ReCaptchaToken,
		}).Post("")
	if captchaErr != nil {
		AbortWithError(c, http.StatusInternalServerError, "Failed retrieving recaptcha response")
		return
	}
	var reCaptcha model.ReCaptchaV3Response
	reErr := json.Unmarshal(reCaptchaResp.Body(), &reCaptcha)
	if reErr != nil {
		AbortWithError(c, http.StatusInternalServerError, "Failed parsing recaptcha response")
		return
	}
	if !reCaptcha.Success {
		AbortWithError(c, http.StatusForbidden, "Failed captcha validation")
		return
	}

	_, claim, err = ca.service.MintForUser(claim, catalog)
	if err != nil {
		if err == model.ErrSignature {
			AbortWithError(c, http.StatusBadRequest, "Invalid signature. Check that you're connected to the right chain.")
			return
		}
		if err == model.ErrExists {
			AbortWithError(c, http.StatusBadRequest, "You've already claimed NFT for this catalog")
			return
		}
		if err == model.ErrKeyword {
			AbortWithError(c, http.StatusBadRequest, "Invalid keywords. Please review the content again")
			return
		}
		AbortWithError(c, http.StatusInternalServerError, "Failed interacting with onchain contract")
		return
	}
	c.JSON(http.StatusOK, claim)
}

// Nft Contract
// @Security     ApiKeyAuth
// @Summary      Nft Contract
// @Description  Gets the current balance of NFT bridge
// @Tags         Nft Bridge
// @Failure      500            {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/bridge/balance [get]
func (nca *ClaimAPI) GetBridgeBalance(c *gin.Context) {
	balance, err := nca.service.GetBalance()
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"balance": balance.String(),
	})
}

// Nft Claim
// @Summary      Nft Claim
// @Description  gets the payload to sign by the user with their wallet
// @Tags         Claiming
// @Param        catalogId  path      string         true  "categoryId"
// @Param        address    path      string         true  "address"
// @Failure      500  {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/claim/{address}/payload/{catalogId} [get]
func (nca *ClaimAPI) SigningPayload(c *gin.Context) {
	catalogId := c.Param("catalogId")
	address := c.Param("address")

	// validate if catalogId exists
	_, catErr := nca.catalogService.GetCatalog(catalogId)
	if catErr != nil {
		if catErr == model.ErrNotFound {
			AbortWithError(c, http.StatusNotFound, "Catalog not found")
			return
		}
		AbortWithError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	// validate if user hasnt already claimed the same category
	_, errClaim := nca.service.GetClaim(catalogId, address)
	if errClaim != model.ErrNotFound {
		// has to be not found
		AbortWithError(c, http.StatusBadRequest, "You have already claimed this NFT")
		return
	}

	// prepare data to sign according to EIP-712
	sd := apitypes.TypedData{}
	copier.Copy(&sd, &model.SignerData) // deep copy object structure (too complex to not many data points required)

	sd.Message["catalogId"] = catalogId
	sd.Message["wallet"] = address
	sd.Domain.Salt = lc.Conf.BlockchainConfig.EIP712TypedData.Salt
	sd.Domain.Name = lc.Conf.BlockchainConfig.EIP712TypedData.Name
	sd.Domain.VerifyingContract = lc.Conf.BlockchainConfig.MailioNFTContractAddress
	sd.Domain.Version = lc.Conf.BlockchainConfig.EIP712TypedData.Version
	sd.Domain.ChainId = math.NewHexOrDecimal256(int64(lc.Conf.BlockchainConfig.DefaultChainId))

	c.JSON(http.StatusOK, sd)
}

// @Summary      Get claimed tx log
// @Description  Reads the claimed transaction log (all mailio claimed NFTs)
// @Tags         Claiming
// @Param        walletaddress  path      string  true   "users wallet address"
// @Param        limit          query     string  false  "limit"
// @Success      200            {array}   model.ClaimPreview
// @Failure      500        {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/user/claims/{walletaddress} [get]
func (nca *ClaimAPI) ListClaimsByUser(c *gin.Context) {
	walletAddress := c.Param("walletaddress")
	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		l, cErr := strconv.Atoi(limitStr)
		if cErr != nil {
			AbortWithError(c, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = l
	}

	claimPreviews, err := nca.service.ReadClaimedTransactionLogs(walletAddress, limit)
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	c.JSON(http.StatusOK, claimPreviews)
}
