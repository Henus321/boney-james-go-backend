package coat

type Coat struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Cost        int32    `json:"cost" bson:"cost"`
	Description string   `json:"description" bson:"description"`
	Model       string   `json:"model" bson:"model"`
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"`
	Colors      []Color  `json:"colors" bson:"colors"`
	Sizes       []string `json:"sizes" bson:"sizes"`
}

type Color struct {
	Label     string   `json:"label" bson:"label"`
	Hex       string   `json:"hex" bson:"hex"`
	PhotoUrls []string `json:"photoUrls" bson:"photoUrls"`
}

type CreateCoatDTO struct {
	Cost        int32    `json:"cost" bson:"cost"`
	Description string   `json:"description" bson:"description"`
	Model       string   `json:"model" bson:"model"`
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"`
	Colors      []Color  `json:"colors" bson:"colors"`
	Sizes       []string `json:"sizes" bson:"sizes"`
}
