package queries

//TODO: use sqlc
const SaveAd = "INSERT INTO ads (created_at, updated_at, city, date, stars, title, address, footage, rooms, floor, specifications, price, premium, ad_id, source, url, building_floors, ad_id_ui) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const SaveLocation = "INSERT INTO locations (id, created_at, updated_at, lat, lng) VALUES (?, ?, ?, ?, ?)"
const SaveRentState = "INSERT INTO rent_daily_city_stats (created_at, average_price, average_price_per_sq, average_footage, city, date, ads_count) VALUES %s"
