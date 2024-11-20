CREATE DATABASE glad
    WITH 
    OWNER = glad_user
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0
    CONNECTION LIMIT = -1;

-- Connect to the new database (only works in psql)
\c glad;

-- Create custom ENUM types (add these before the table creations)
CREATE TYPE course_status AS ENUM ('draft'
    , 'archived'
    , 'open' 
    , 'expense-submitted'
    , 'expense-declined'
    , 'closed'
    , 'active'
    , 'declined'
    , 'submitted'
    , 'canceled'
    , 'inactive');
CREATE TYPE course_type AS ENUM ('in-person'
    , 'online');
CREATE TYPE timezone_type AS ENUM ('EST'
    , 'CST'
    , 'MST'
    , 'PST');
CREATE TYPE account_type AS ENUM ('teacher'
    , 'assistant-teacher'
    , 'organizer'
    , 'member'
    , 'user');
CREATE TYPE center_mode AS ENUM ('in-person'
    , 'online');

-- Create tables
CREATE TABLE IF NOT EXISTS tenant (
    id BIGSERIAL PRIMARY KEY,
    -- TODO: Name need not be unique, but name and country together must be unique
    name VARCHAR(255) NOT NULL UNIQUE,
    country VARCHAR(128) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_tenant_name ON tenant(name);

-- CENTER entity
CREATE TABLE IF NOT EXISTS center (
    id SERIAL PRIMARY KEY,
    -- Note: ext_id is salesforce id
    ext_id VARCHAR(32) NOT NULL UNIQUE,
    -- TODO: tenant_id BIGINT NOT NULL REFERENCES tenant(id), -- ON DELETE RESTRICT, 
    -- (Note: Do not want to delete tenant if center exists; maybe we can set to NULL as an alternate choice)

    -- Note: 'name' is needed in SalesForce. Set as 'L-<id>'
    center_name VARCHAR(80) NOT NULL,
    
    -- Note: location format: {"street_1": ..., "street_2": ..., "city": ..., "state": ..., "zip": ..., "country": ...}
    -- When multitenancy is introduced, then country can be removed.
    location JSONB,

    -- Note: geo_location format: {"lat": ..., "long": ...}
    geo_location JSONB,
    -- maximum occupancy
    capacity INTEGER,
    mode center_mode DEFAULT 'in-person',
    webpage VARCHAR(255),
    is_national_center BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_center_ext_id ON center(ext_id);
CREATE INDEX idx_center_center_name ON center(center_name);

CREATE TABLE IF NOT EXISTS center_contact (
    center_id INT NOT NULL REFERENCES center(id) ON DELETE CASCADE,
    name VARCHAR(255),
    phone VARCHAR(32),
    email VARCHAR(80),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_center_contact_center_id ON center_contact(center_id);

-- COURSE entity
CREATE TABLE IF NOT EXISTS course (
    id BIGSERIAL PRIMARY KEY,
    -- Note: ext_id is salesforce id
    -- When course is created outside of salesforce, ext_id will be NULL
    ext_id VARCHAR(32) UNIQUE,
    -- TODO: What's CType ID? How is it used?
    -- TODO: tenant_id BIGINT NOT NULL REFERENCES tenant(id), -- ON DELETE RESTRICT, 
    -- (Note: Do not want to delete tenant if course exists; maybe we can set to NULL as an alternate choice)
    name VARCHAR(128) NOT NULL,
    notes VARCHAR(1024), -- TODO: check the size of this column
    status course_status NOT NULL DEFAULT 'draft',
    max_attendees INTEGER,
    timezone timezone_type,
    -- Note: location format: {"street_1": ..., "street_2": ..., "city": ..., "state": ..., "zip": ..., "country": ...}
    -- When multitenancy is introduced, then country can be removed.
    location JSONB,
    center_id BIGINT NOT NULL REFERENCES center(id) ON DELETE RESTRICT,
    ctype course_type NOT NULL DEFAULT 'in-person',
    num_attendees INTEGER DEFAULT 0,
    is_auto_approve BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- CREATE INDEX idx_course_tenant_id ON course(tenant_id);
CREATE INDEX idx_course_ext_id ON course(ext_id);

-- ACCOUNT entity
CREATE TABLE IF NOT EXISTS account (
    id BIGSERIAL PRIMARY KEY,
    -- Note: ext_id is salesforce id
    ext_id VARCHAR(32) NOT NULL UNIQUE,
    -- TODO: tenant_id BIGINT NOT NULL REFERENCES tenant(id), -- ON DELETE RESTRICT, 
    -- (Note: Do not want to delete tenant if course exists; maybe we can set to NULL as an alternate choice)

    -- Note: username is to link the account with the logged in user
    username VARCHAR(128) NOT NULL,
    first_name VARCHAR(40),
    last_name VARCHAR(80),
    phone VARCHAR(32),
    email VARCHAR(80),
    type account_type,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    -- Q: Is it unique? How about kids' account?
    -- UNIQUE(tenant_id, username)
);
-- CREATE INDEX idx_account_tenant_id ON account(tenant_id);
CREATE INDEX idx_account_ext_id ON account(ext_id);
CREATE INDEX idx_account_username ON account(username);
CREATE INDEX idx_account_type ON account(type);
-- TODO: Need indexes for email and phone?

-- Notes: Max 3 organizers per course
-- Notes: Organizer is not mandatory for a course (Confirm)
CREATE TABLE IF NOT EXISTS course_organizer (
    course_id BIGINT NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    organizer_id BIGINT NOT NULL REFERENCES account(id) ON DELETE RESTRICT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_course_organizer_course_id ON course_organizer(course_id);
CREATE INDEX idx_course_organizer_organizer_id ON course_organizer(organizer_id);

-- Notes: Max 1 contact per course
-- Notes: Contact is not mandatory for a course
CREATE TABLE IF NOT EXISTS course_contact (
    course_id BIGINT NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    contact_id BIGINT NOT NULL REFERENCES account(id) ON DELETE RESTRICT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_course_contact_course_id ON course_contact(course_id);
CREATE INDEX idx_course_contact_contact_id ON course_contact(contact_id);

-- Notes: Max 3 teachers per course
-- Notes: As per SF data model, even teachers are optional for a course
CREATE TABLE IF NOT EXISTS course_teacher (
    course_id BIGINT NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    teacher_id BIGINT NOT NULL REFERENCES account(id) ON DELETE RESTRICT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_course_teacher_course_id ON course_teacher(course_id);
CREATE INDEX idx_course_teacher_teacher_id ON course_teacher(teacher_id);

CREATE TABLE IF NOT EXISTS course_timing (
    course_id BIGINT NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    -- Note: ext_id is salesforce id
    ext_id VARCHAR(32) NOT NULL UNIQUE,
    -- Note: 'name' is needed in SalesForce. Set as 'D-mmddyyyy'
    -- Note: 'timezone' is needed in SalesForce. Set the value from 'course' table
    course_date DATE,
    start_time TIME,
    end_time TIME,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_course_timings_course_id ON course_timing(course_id);
CREATE INDEX idx_course_timing_ext_id ON course_timing(ext_id);

-- Notes: Notify is not mandatory for a course
CREATE TABLE IF NOT EXISTS course_notify (
    course_id BIGINT NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    notify_id BIGINT NOT NULL REFERENCES account(id) ON DELETE RESTRICT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_course_notify_course_id ON course_notify(course_id);
-- CREATE INDEX idx_course_notify_notify_id ON course_notify(notify_id); -- TODO: May be this is not required
