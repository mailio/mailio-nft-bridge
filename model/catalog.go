package model

const CatalogTable = "catalog"

// Catalog serves as knowledge catalog high level description
type Catalog struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name" validate:"required,min=3,max=255"`
	Type          string `json:"type" validate:"required" oneof:"video,article,podcast,podcast-episode,virtual-event"`
	Description   string `json:"description" validate:"required,min=3,max=1000"`
	ContentLink   string `json:"contentLink" validate:"required,min=3,max=2000"`
	Keywords      string `json:"keywords" validate:"required,min=3,max=1000"` // comma separated list of keywords
	VideoLink     string `json:"videoLink,omitempty"`                         // YouTube or similar link
	ImageLink     string `json:"imageLink,omitempty"`                         // CID/hash of the image
	NftTokensUsed int    `json:"nftTokensUsed"`                               //currently minted tokens for the catalog
	Modified      int64  `json:"modified"`
	Created       int64  `json:"created"`
}
