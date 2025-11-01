-- Add inputfileids column to projects table
-- This column stores JSON array of input file IDs, similar to cfileids

USE scvis_local;

ALTER TABLE projects ADD COLUMN inputfileids JSON NOT NULL DEFAULT (JSON_ARRAY());
