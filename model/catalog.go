package model

const CatalogTable = "catalog"

// Catalog serves as knowledge catalog high level description
type Catalog struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name" validate:"required,min=3,max=255"`
	Description    string `json:"description"`
	VideoLink      string `json:"videoLink,omitempty"`                              // YouTube or similar link
	ImageLink      string `json:"imageLink,omitempty"`                              // NFT Representative image (IPFS preffered)
	NftTokensTotal int    `json:"nftTokensTotal" validate:"required,numeric,min=1"` // max number of minted tokens for the catalog
	NftTokensUsed  int    `json:"nftTokensUsed"`                                    //currently minted tokens for the catalog
	Modified       int64  `json:"modified"`
	Created        int64  `json:"created"`
}
