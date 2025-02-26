-- building
-- id: Primary key, integer, auto-increment
-- name: String, unique
-- address: Text

CREATE TABLE IF NOT EXISTS building (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    address TEXT
);

-- Apartment
-- id: Primary key, integer, auto-increment
-- building_id: Foreign key, integer (references building.id)
-- number: String
-- floor: Integer
-- sq_meters: Integer

CREATE TABLE IF NOT EXISTS apartment (
    id SERIAL PRIMARY KEY,
    building_id INTEGER,
    number TEXT,
    floor INTEGER,
    sq_meters INTEGER,
    FOREIGN KEY (building_id) REFERENCES building(id)
);
