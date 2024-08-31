DROP TABLE IF EXISTS coat CASCADE;

CREATE TABLE public.coat (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        model VARCHAR(100) NOT NULL UNIQUE,
        name VARCHAR(255) NOT NULL,
        description VARCHAR(255) NOT NULL
);

CREATE TABLE public.coat_option (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        colorLabel VARCHAR(100) NOT NULL,
        colorHex VARCHAR(25) NOT NULL,
        cost INTEGER NOT NULL,
        sizes VARCHAR(25)[],
        photoUrls VARCHAR(255)[],
        coatId UUID REFERENCES coat(id)
);

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

