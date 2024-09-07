package shop

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/jackc/pgtype"
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
