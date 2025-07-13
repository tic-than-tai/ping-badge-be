-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    profile_picture_url VARCHAR(255),
    bio TEXT,
    role VARCHAR(20) DEFAULT 'USER' CHECK (role IN ('USER', 'ORGANIZER')),
    privacy_setting VARCHAR(20) DEFAULT 'public',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Organizations table
CREATE TABLE organizations (
    org_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_name VARCHAR(255) UNIQUE NOT NULL,
    org_email VARCHAR(100) UNIQUE NOT NULL,
    org_logo_url VARCHAR(255),
    user_id_owner UUID REFERENCES users(user_id),
    description TEXT,
    website_url VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Organization admins table
CREATE TABLE organization_admins (
    admin_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(org_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);

-- Badges table
CREATE TABLE badges (
    badge_def_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(org_id) ON DELETE CASCADE,
    badge_name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255) NOT NULL,
    criteria TEXT,
    badge_type VARCHAR(20) NOT NULL DEFAULT 'instant' CHECK (badge_type IN ('instant', 'cumulative')),
    rule_config JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(org_id, badge_name)
);

-- Issued badges table
CREATE TABLE issued_badges (
    issued_badge_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    badge_def_id UUID NOT NULL REFERENCES badges(badge_def_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    org_id UUID NOT NULL REFERENCES organizations(org_id) ON DELETE CASCADE,
    issue_date TIMESTAMP DEFAULT NOW(),
    verification_code VARCHAR(255) UNIQUE NOT NULL,
    source_type VARCHAR(50),
    source_id UUID,
    cumulative_progress_at_issuance NUMERIC,
    cumulative_unit VARCHAR(50),
    additional_data JSONB,
    status VARCHAR(20) DEFAULT 'issued',
    blockchain_tx_id VARCHAR(255)
);

-- User badge progress table
CREATE TABLE user_badge_progress (
    progress_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    badge_def_id UUID NOT NULL REFERENCES badges(badge_def_id) ON DELETE CASCADE,
    progress_value NUMERIC DEFAULT 0,
    unit VARCHAR(50),
    is_qualified BOOLEAN DEFAULT FALSE,
    last_updated TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, badge_def_id)
);

-- Activities table
CREATE TABLE activities (
    activity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(org_id) ON DELETE CASCADE,
    activity_name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    location VARCHAR(255),
    badge_def_id UUID REFERENCES badges(badge_def_id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Activity participation table
CREATE TABLE activity_participation (
    participation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_id UUID NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'registered',
    proof_of_participation_url VARCHAR(255),
    issued_badge_id UUID UNIQUE REFERENCES issued_badges(issued_badge_id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Badge views table
CREATE TABLE badge_views (
    view_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issued_badge_id UUID NOT NULL REFERENCES issued_badges(issued_badge_id) ON DELETE CASCADE,
    viewer_ip_address VARCHAR(45),
    view_timestamp TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_issued_badges_user_badge ON issued_badges(user_id, badge_def_id);
CREATE INDEX idx_issued_badges_org_badge ON issued_badges(org_id, badge_def_id);
CREATE INDEX idx_user_badge_progress_user_badge ON user_badge_progress(user_id, badge_def_id);
CREATE INDEX idx_activities_org_id ON activities(org_id);
CREATE INDEX idx_activity_participation_activity_user ON activity_participation(activity_id, user_id);
