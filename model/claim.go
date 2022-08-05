package model

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const ClaimTable = "claim"
const ClaimFingerprintTable = "fingerprint"

type Claim struct {
	CatalogId      string         `json:"catalogId" validate:"required"`      // categoryId to be claimed
	WalletAddress  string         `json:"walletAddress" validate:"required"`  // publickey of the user retrieved from wallet
	MailioAddress  string         `json:"mailioAddress,omitempty"`            // optional mailio address
	Signature      string         `json:"signature" validate:"required"`      // signature of categoryId + nonce
	ReCaptchaToken string         `json:"recaptchaToken" validate:"required"` // recaptcha v3 token // required
	GasPrice       uint64         `json:"gasPrice"`                           // gas price of the transaction
	TxHash         string         `json:"txHash,omitempty"`                   // transaction hash of the transaction
	TokenUri       string         `json:"tokenUri,omitempty"`                 // token uri
	VisitorId      string         `json:"visitorId" validate:"required"`      // visitor id
	Keywords       []ClaimKeyword `json:"keywords,omitempty"`                 // keywords (not need to be stored in db)
	Created        int64          `json:"created"`
}

// fingerprinting each catalogId claim in order to prevent users
// getting the same catalogId claim multiple times
type ClaimFingerprint struct {
	CatalogId string `json:"catalogId" validate:"required"` // categoryId to be claimed
	VisitorId string `json:"visitorId" validate:"required"` // fingerprint of the user
}

// preview of the claimed token (not need to be stored in db)
// tokenId is retrieved from the blockchain
type ClaimPreview struct {
	Claim
	TokenId  uint64 `json:"tokenId"`
	TxStatus uint64 `json:"txStatus"` // 1 = success, 0 = fail
}

type ClaimKeyword struct {
	Word string `json:"word"`
}

type ReCaptchaV3Response struct {
	Success            bool     `json:"success"`
	ChallengeTimestamp string   `json:"challenge_ts"` //  timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname           string   `json:"hostname"`     // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes         []string `json:"error-codes"`  // optional error codes
}

// EIP-712 -- https://eips.ethereum.org/EIPS/eip-712
var (
	SignerData = apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:    "",
			Version: "",
			ChainId: math.NewHexOrDecimal256(0),
		},
		Message: map[string]interface{}{
			"catalogId": "string",
			"wallet":    "address",
		},
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"claim": {
				{Name: "catalogId", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
		},
		PrimaryType: "claim",
	}
)
