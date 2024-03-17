CREATE TABLE ads (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    city TEXT NULL,
    date TIMESTAMP NULL,
    stars INTEGER NULL,
    title TEXT NULL,
    address TEXT NULL,
    footage NUMERIC NULL,
    rooms INTEGER NULL,
    floor INTEGER NULL,
    specifications TEXT NULL,
    price INTEGER NULL,
    premium BOOLEAN NULL,
    ad_id TEXT NULL,
    source TEXT NULL,
    url TEXT NULL,
    building_floors INTEGER NULL,
    ad_id_ui INTEGER NULL
);

CREATE INDEX IF NOT EXISTS idx_ads_deleted_at ON ads (deleted_at);

