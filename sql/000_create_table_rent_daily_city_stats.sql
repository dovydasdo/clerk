CREATE TABLE rent_daily_city_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATE NULL,
    updated_at DATE NULL,
    deleted_at DATE NULL,
    average_price INTEGER NULL,
    average_price_per_sq REAL NULL,
    average_footage REAL NULL,
    city TEXT NULL,
    date DATE NULL,
    ads_count INTEGER NULL
    source TEXT NULL,
);
