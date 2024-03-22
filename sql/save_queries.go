package queries

//TODO: use sqlc
const SaveAd = "INSERT INTO ads (created_at, updated_at, city, date, stars, title, address, footage, rooms, floor, specifications, price, premium, ad_id, source, url, building_floors, ad_id_ui) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)"
const SaveLocation = "INSERT INTO locations (id, created_at, updated_at, lat, lng) VALUES ($1, $2, $3, $4, $5)"
