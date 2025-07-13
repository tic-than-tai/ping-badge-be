-- Drop tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS badge_views;
DROP TABLE IF EXISTS activity_participation;
DROP TABLE IF EXISTS activities;
DROP TABLE IF EXISTS user_badge_progress;
DROP TABLE IF EXISTS issued_badges;
DROP TABLE IF EXISTS badges;
DROP TABLE IF EXISTS organization_admins;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
