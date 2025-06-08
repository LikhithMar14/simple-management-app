-- +goose Up
-- Drop existing foreign key constraint (if exists)
ALTER TABLE IF EXISTS cars DROP CONSTRAINT IF EXISTS fk_engine_id;

-- Create engines table first
CREATE TABLE IF NOT EXISTS engines (
    id UUID PRIMARY KEY,
    displacement INT NOT NULL,
    number_of_cylinders INT NOT NULL,
    car_range INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create cars table
CREATE TABLE IF NOT EXISTS cars (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    year VARCHAR(4) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    fuel_type VARCHAR(50) NOT NULL,
    engine_id UUID NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Now truncate tables to clear existing data
TRUNCATE TABLE cars;
TRUNCATE TABLE engines;

-- Add foreign key constraint on engine_id in cars table
ALTER TABLE cars ADD CONSTRAINT fk_engine_id FOREIGN KEY (engine_id) REFERENCES engines(id) ON DELETE CASCADE;

-- Insert dummy data into the engines table
INSERT INTO engines (id, displacement, number_of_cylinders, car_range) VALUES
    ('e1f86b1a-0873-4c19-bae2-fc60329d0140', 2000, 4, 600),
    ('f4a9c66b-8e38-419b-93c4-215d5cefb318', 1600, 4, 550),
    ('cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 3000, 6, 700),
    ('9746be12-07b7-42a3-b8ab-7d1f209b63d7', 1800, 4, 500);

-- Insert dummy data into the cars table
INSERT INTO cars (id, name, year, brand, fuel_type, engine_id, price) VALUES
    ('c7c1a6d5-1ec4-4c64-a59a-8a2f6f3d2bf3', 'Honda Civic', '2023', 'Honda', 'Gasoline', 'e1f86b1a-0873-4c19-bae2-fc60329d0140', 25000.00),
    ('9d6a56f8-79c3-4931-a5c0-6b290c84ba2f', 'Toyota Corolla', '2022', 'Toyota', 'Gasoline', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', 22000.00),
    ('9b9437c4-3ed1-45a5-b240-0fe3e24e0e4e', 'Ford Mustang', '2024', 'Ford', 'Gasoline', 'cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 40000.00),
    ('5e9df51a-8d7a-4d84-9c58-4ccfe5c7db06', 'BMW 3 Series', '2023', 'BMW', 'Gasoline', '9746be12-07b7-42a3-b8ab-7d1f209b63d7', 35000.00);

-- +goose Down
ALTER TABLE cars DROP CONSTRAINT IF EXISTS fk_engine_id;
DROP TABLE IF EXISTS cars;
DROP TABLE IF EXISTS engines;