package coat

import "github.com/gofrs/uuid"

type Coat struct {
	ID          string `json:"id"`
	Model       string `json:"model"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CoatOption struct {
	ID         uuid.UUID `json:"id"`
	ColorLabel string    `json:"colorLabel"`
	ColorHex   string    `json:"colorHex"`
	Cost       int32     `json:"cost"`
	Sizes      []string  `json:"sizes"`
	PhotoUrls  []string  `json:"photoUrls"`
	CoatID     uuid.UUID `json:"coatId"`
}

type CoatWithOption struct {
	ID          uuid.UUID    `json:"id"`
	Model       string       `json:"model"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CoatOptions []CoatOption `json:"coatOptions"`
}

type CreateCoatInput struct {
	Model       string `json:"model"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCoatOptionInput struct {
	CoatID     string   `json:"coatId"`
	ColorLabel string   `json:"colorLabel"`
	ColorHex   string   `json:"colorHex"`
	Cost       int32    `json:"cost"`
	Sizes      []string `json:"sizes"`
	PhotoUrls  []string `json:"photoUrls"`
}
