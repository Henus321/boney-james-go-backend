package shop

import (
	"github.com/gofrs/uuid"
)

type Shop struct {
	ID         uuid.UUID `json:"id"`
	CityID     uuid.UUID `json:"cityId"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Street     string    `json:"street"`
	Subway     string    `json:"subway"`
	OpenPeriod string    `json:"OpenPeriod"`
}

type ShopCity struct {
	ID        uuid.UUID `json:"id"`
	CityName  string    `json:"cityName"`
	CityLabel string    `json:"cityLabel"`
}

type ShopType struct {
	ID        uuid.UUID `json:"id"`
	TypeName  string    `json:"typeName"`
	TypeLabel string    `json:"typeLabel"`
}

type ShopWithType struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	Street     string     `json:"street"`
	Subway     string     `json:"subway"`
	OpenPeriod string     `json:"OpenPeriod"`
	ShopCity   ShopCity   `json:"shopCity"`
	ShopTypes  []ShopType `json:"shopTypes"`
}
