package coat

type Coat struct {
	ID          string `json:"id"`
	Model       string `json:"model"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CoatOption struct {
	ID         string   `json:"id"`
	ColorLabel string   `json:"colorLabel"`
	ColorHex   string   `json:"colorHex"`
	Cost       int32    `json:"cost"`
	Sizes      []string `json:"sizes"`
	PhotoUrls  []string `json:"photoUrls"`
	CoatId     string   `json:"coatId"`
}

type CoatWithOption struct {
	ID          string       `json:"id"`
	Model       string       `json:"model"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CoatOptions []CoatOption `json:"coatOptions"`
}

type CreateCoatDTO struct {
	Model       string `json:"model"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
