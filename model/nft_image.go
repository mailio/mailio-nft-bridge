package model

type NftImageUploadResponse struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Size string `json:"size"`
}

type NftImage struct {
	PinLsList   NftKeys        `json:"Keys"`
	PinLsObject NftPinLsObject `json:"PinLsObject"`
}

type NftKeys map[string]interface{}

type NftPinLsObject struct {
	Cid  string `json:"Cid"`
	Type string `json:"Type"`
}

type NftPins struct {
	Pins []string `json:"pins"`
}

type Erc721Json struct {
	// Title      string            `json:"title"`
	// Type       string            `json:"type"`
	// Properties *Erc721Properties `json:"properties,omitempty"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Image           string             `json:"image"`
	YoutubeUrl      string             `json:"youtube_url,omitempty"`
	BackgroundColor string             `json:"background_color,omitempty"`
	ExternalUrl     string             `json:"external_url,omitempty"`
	Attributes      []*Erc721Attribute `json:"attributes"`
}

type Erc721Attribute map[string]interface{}

// type Erc721Properties struct {
// 	Name        *Erc721Property `json:"name,omitempty"`
// 	Description *Erc721Property `json:"description,omitempty"`
// 	Image       *Erc721Property `json:"image,omitempty"`
// }

// type Erc721Property struct {
// 	Type        string `json:"type"`
// 	Description string `json:"description,omitempty"`
// }
