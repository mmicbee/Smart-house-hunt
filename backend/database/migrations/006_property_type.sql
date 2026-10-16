ALTER TABLE properties ADD COLUMN property_type TEXT NOT NULL DEFAULT 'apartment';
CREATE INDEX IF NOT EXISTS idx_properties_type ON properties(property_type);
