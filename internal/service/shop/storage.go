package shop

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"log"
)

type Storage struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewStorage(client postgresql.Client, logger *logging.Logger) *Storage {
	return &Storage{
		client: client,
		logger: logger,
	}
}

func (s *Storage) GetAllShops(ctx context.Context, cityId *string, typeId *string) (*[]ShopWithType, error) {
	const op = "coat.storage.GetAllShops"

	query := `SELECT
			id,
			name,
			phone,
			street,
			subway,
			openPeriod,
			sp.cityId,
			cityName,
			cityLabel,
			st.typeId,
			typeName,
			typeLabel
	FROM shop_with_type as swp
    INNER JOIN
        (SELECT sh.id, sh.cityId, sh.name,sh.phone, sh.street,sh.subway, sh.openPeriod ,shct.cityName, shct.cityLabel, shct.id as shopCityId
            FROM shop as sh LEFT JOIN shop_city as shct ON cityId = shct.id) as sp
                ON sp.id = shopId
    INNER JOIN (SELECT id as typeId, typeName, typeLabel FROM shop_type) as st ON swp.shopTypeId = st.typeId`
	// TODO WHERE (cityId = $1 or $1 IS NULL) AND (typeId = $2 or $2 IS NULL)`
	// 1 кейс - Передаю пустоту - показывай все
	// 2 кейс - Фильтрация по typeId должна убирать магазины без таких типов в листе, но сейчас убирает типы из самого листа внутри магазина
	// 3 кейс - Нужно прислать все уникальные типы магазинов и городов для селектов, сделать это отдельной функцией или запихнуть сюда?
	log.Printf("cityId: %v, typeId %v", cityId, typeId)
	//rows, err := s.client.Query(ctx, query, cityId, typeId)

	rows, err := s.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get shops: %w", op, err)
	}
	defer rows.Close()

	var swts []ShopWithType

	for rows.Next() {
		var (
			id         pgtype.UUID
			name       pgtype.Text
			phone      pgtype.Text
			street     pgtype.Text
			subway     pgtype.Text
			openPeriod pgtype.Text
			cityId     pgtype.UUID
			cityName   pgtype.Text
			cityLabel  pgtype.Text
			typeId     pgtype.UUID
			typeName   pgtype.Text
			typeLabel  pgtype.Text
		)

		err = rows.Scan(
			&id,
			&name,
			&phone,
			&street,
			&subway,
			&openPeriod,
			&cityId,
			&cityName,
			&cityLabel,
			&typeId,
			&typeName,
			&typeLabel,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get shop: %w", op, err)
		}

		index := IndexOfShop(swts, id.Bytes)

		if id.Status == pgtype.Present &&
			name.Status == pgtype.Present &&
			phone.Status == pgtype.Present &&
			street.Status == pgtype.Present &&
			subway.Status == pgtype.Present &&
			openPeriod.Status == pgtype.Present &&
			cityId.Status == pgtype.Present &&
			cityLabel.Status == pgtype.Present &&
			cityName.Status == pgtype.Present &&
			typeId.Status == pgtype.Present &&
			typeName.Status == pgtype.Present &&
			typeLabel.Status == pgtype.Present {
			if index == -1 {
				swts = append(swts, ShopWithType{
					ID:         id.Bytes,
					Name:       name.String,
					Phone:      phone.String,
					Street:     street.String,
					Subway:     subway.String,
					OpenPeriod: openPeriod.String,
					ShopCity: ShopCity{
						ID:        cityId.Bytes,
						CityName:  cityName.String,
						CityLabel: cityLabel.String,
					},
					ShopTypes: []ShopType{{
						ID:        typeId.Bytes,
						TypeName:  typeName.String,
						TypeLabel: typeLabel.String,
					}},
				})
			} else {
				swts[index].ShopTypes = append(swts[index].ShopTypes, ShopType{
					ID:        typeId.Bytes,
					TypeName:  typeName.String,
					TypeLabel: typeLabel.String,
				})
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to get all shops: %w", op, err)
	}

	return &swts, nil
}

func (s *Storage) GetShopByID(ctx context.Context, id string) (*ShopWithType, error) {
	const op = "coat.storage.GetShopByID"

	query := `SELECT 
			id,
			name,
			phone,
			street,
			subway,
			openPeriod,
			sp.cityId,
			cityName,
			cityLabel,
			st.typeId,
			typeName,
			typeLabel
		FROM shop_with_type as swp
		INNER JOIN
			(SELECT sh.id, sh.cityId, sh.name,sh.phone, sh.street,sh.subway, sh.openPeriod ,shct.cityName, shct.cityLabel, shct.id as shopCityId
		FROM shop as sh LEFT JOIN shop_city as shct ON cityId = shct.id) as sp
			ON sp.id = shopId
		INNER JOIN (SELECT id as typeId, typeName, typeLabel FROM shop_type) as st ON swp.shopTypeId = st.typeId
		WHERE swp.shopId = $1`

	rows, err := s.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get coats: %w", op, err)
	}
	defer rows.Close()

	var swt ShopWithType

	for rows.Next() {
		var (
			typeId    pgtype.UUID
			typeName  pgtype.Text
			typeLabel pgtype.Text
		)

		err = rows.Scan(
			&swt.ID,
			&swt.Name,
			&swt.Phone,
			&swt.Street,
			&swt.Subway,
			&swt.OpenPeriod,
			&swt.ShopCity.ID,
			&swt.ShopCity.CityName,
			&swt.ShopCity.CityLabel,
			&typeId,
			&typeName,
			&typeLabel,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get shop: %w", op, err)
		}

		if typeId.Status == pgtype.Present &&
			typeName.Status == pgtype.Present &&
			typeLabel.Status == pgtype.Present {
			swt.ShopTypes = append(swt.ShopTypes, ShopType{
				ID:        typeId.Bytes,
				TypeName:  typeName.String,
				TypeLabel: typeLabel.String,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to get coat by id: %w", op, err)
	}

	return &swt, nil
}

// IndexOfShop ??? Куда такие хелперы с конкретным значением. GO generic ok?
func IndexOfShop(shops []ShopWithType, shopId uuid.UUID) int {
	for i, v := range shops {
		if v.ID == shopId {
			return i
		}
	}
	return -1
}
