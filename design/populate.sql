CREATE TABLE phone_purchases(
    model_name VARCHAR,
    client_name VARCHAR,
    price REAL,
    purchase_date DATETIME,
    phisical_store VARCHAR,
    selling_employer VARCHAR
);

INSERT INTO phone_purchases VALUES(
    'iphone 4s',
    'Emanuel',
    456.0,
    '2025-01-02',
    'St. Some random place 423, Here',
    'John Doe'
);

INSERT INTO phone_purchases VALUES(
    'samsung s24',
    'Pinco Pallo',
    950.0,
    '2022-01-02',
    'St. NewOne, There',
    'Johnna Doe'
);

CREATE TABLE gym_subscriptions(
    gym_name VARCHAR,
    client_name VARCHAR,
    subscription_price REAL,
    subscription_start_date DATETIME,
    subscription_end_date DATETIME,
    purchase_date DATETIME,
    gym_address VARCHAR,
    selling_employer VARCHAR
);

INSERT INTO gym_subscriptions VALUES(
    'TrainWithUs',
    'Pinco Pallo',
    650.0,
    '2025-01-01',
    '2025-09-01',
    '2025-01-01',
    'St. Where the gym is, Space',
    'MyFavourite Seller'
);

INSERT INTO gym_subscriptions VALUES(
    'TrainHard',
    'Simon',
    250.0,
    '2026-01-01',
    '2026-05-01',
    '2026-01-01',
    'St. Unknown, NewLand',
    'BadSeller'
);