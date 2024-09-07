----- COAT -----
DROP TABLE IF EXISTS coat CASCADE;

CREATE TABLE public.coat (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        model VARCHAR(100) NOT NULL UNIQUE,
        name VARCHAR(255) NOT NULL,
        description VARCHAR(255) NOT NULL
);

CREATE TABLE public.coat_option (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        coatId UUID,
        colorLabel VARCHAR(100) NOT NULL,
        colorHex VARCHAR(25) NOT NULL,
        cost INTEGER NOT NULL,
        sizes VARCHAR(25)[],
        photoUrls VARCHAR(255)[],

        FOREIGN KEY (coatId) REFERENCES public.coat (id) ON DELETE CASCADE
);

SELECT * FROM coat
                  LEFT JOIN coat_option
                            ON coat.id = coat_option.coatid;

INSERT INTO coat (
    model,
    name,
    description
) VALUES (
    'A-001',
    'Короткое пальто',
    'Элегантное короткое пальто на кнопках, прямой крой и воротник-стойка'
);

INSERT INTO coat_option (
            colorLabel,
            colorHex,
            cost,
            sizes,
            photoUrls,
            coatId
) VALUES (
            'Бирюзовый',
            '#2e9ec0',
            10999,
            ARRAY ['XS 42', 'S 44', 'M 46', 'L 48'],
            ARRAY [
            'https://firebasestorage.googleapis.com/v0/b/boney-james-c978c.appspot.com/o/2016%2F2021-a-005-6.jpg?alt=media&token=669b7b47-4f74-4684-a50b-8415e3e75cde',
            'https://firebasestorage.googleapis.com/v0/b/boney-james-c978c.appspot.com/o/2016%2F2021-a-005-7.jpg?alt=media&token=7609c089-9e0b-44a8-afa8-1988548d7acf',
            'https://firebasestorage.googleapis.com/v0/b/boney-james-c978c.appspot.com/o/2016%2F2021-a-005-8.jpg?alt=media&token=d8dcd203-61fe-4177-baa1-df2ad0ba2f7c'
            ],
            '0f656773-585b-403d-abde-12117ec860d0'
);

INSERT INTO coat_option (
    colorLabel,
    colorHex,
    cost,
    sizes,
    photoUrls,
    coatId
) VALUES (
             'Темно-синий',
             '#080e61',
             11999,
             ARRAY ['XS 42', 'S 44', 'M 46', 'L 48'],
             ARRAY [
                 'https://firebasestorage.googleapis.com/v0/b/boney-james-c978c.appspot.com/o/2016%2F2021-a-005-1.jpg?alt=media&token=4f301a1d-77d8-4e15-9117-a40ae1107a42',
                 'https://firebasestorage.googleapis.com/v0/b/boney-james-c978c.appspot.com/o/2016%2F2021-a-005-2.jpg?alt=media&token=3c6cce45-978c-4b8f-a25f-7529c36d948d',
                 'https://firebasestorage.googleapis.com/v0/b/boney-james-c978c.appspot.com/o/2016%2F2021-a-005-3.jpg?alt=media&token=698d4f96-c5b4-4bcf-91e5-71cbb3147290'
                 ],
             '0f656773-585b-403d-abde-12117ec860d0'
 );

SELECT * FROM coat LEFT JOIN public.coat_option ON coat.id = coat_option.coatId;

----- SHOP -----
DROP TABLE IF EXISTS shop CASCADE;

CREATE TABLE public.shop (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        cityId UUID NOT NULL,
        name  VARCHAR(255) NOT NULL,
        phone  VARCHAR(25) NOT NULL,
        street  VARCHAR(255) NOT NULL,
        subway  VARCHAR(255) NOT NULL,
        openPeriod  VARCHAR(100) NOT NULL,

        FOREIGN KEY (cityId) REFERENCES public.shop_city (id)
);

CREATE TABLE public.shop_city (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        cityName  VARCHAR(255) NOT NULL,
        cityLabel  VARCHAR(255) NOT NULL
);

CREATE TABLE public.shop_type (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        typeName  VARCHAR(255) NOT NULL,
        typeLabel  VARCHAR(255) NOT NULL
);

CREATE TABLE public.shop_with_type (
        shopId UUID,
        shopTypeId UUID,

        PRIMARY KEY (shopId, shopTypeId),
        FOREIGN KEY (shopId) REFERENCES public.shop (id),
        FOREIGN KEY (shopTypeId) REFERENCES public.shop_type (id)
);

INSERT INTO shop (
        cityId,
        name,
        phone,
        street,
        subway,
        openPeriod
) VALUES (
        'fdb2de16-75b1-4ba4-9702-d90ca05bd110',
        'ТЦ Европейский',
        '+7 (495) 921-34-44',
        'ул.Площадь Киевского Вокзала',
        'Метро Киевская',
        '10:00-22:00'
);

INSERT INTO shop_city (
        cityName,
        cityLabel
) VALUES (
        'moscow',
        'Москва'
);

INSERT INTO shop_type (
        typeName,
        typeLabel
) VALUES (
        'female',
        'Одежда для женщин'
);

INSERT INTO shop_with_type (
        shopId,
        shopTypeId
) VALUES (
        '02cf3cf3-67f2-4268-a6cb-f117c6517085',
        'd3afe521-92e1-4c97-b556-201169520b12'
);

SELECT
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
    INNER JOIN (SELECT id as typeId, typeName, typeLabel FROM shop_type) as st ON swp.shopTypeId = st.typeId;

SELECT
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
    WHERE swp.shopId = '02cf3cf3-67f2-4268-a6cb-f117c6517085';