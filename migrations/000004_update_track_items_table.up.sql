-- Migration to add subtype and inbound_rule to track_items table
ALTER TABLE track_items ADD COLUMN subtype TEXT;
ALTER TABLE track_items ADD COLUMN inbound_rule TEXT;
