-- +goose Up
-- +goose StatementBegin

CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    creation_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type CHAR(2) NOT NULL, -- GG, BG, CE (Co-ed), UN (Unknown)
    debut_date DATE,
    company_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,
    idol_id INTEGER NOT NULL,
    since_date DATE,
    until_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE idols (
    id SERIAL PRIMARY KEY,
    stage_name VARCHAR(255) NOT NULL, -- Primary name used in performances
    name VARCHAR(255), -- Legal name / birth name
    gender CHAR NOT NULL DEFAULT 'U', -- M, F, O (Other), U (Unknown)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE idol_info (
    id SERIAL PRIMARY KEY,
    idol_id INTEGER NOT NULL,
    birthdate DATE,
    height_cm FLOAT,
    weight_kg FLOAT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Foreign key constraints
ALTER TABLE groups ADD CONSTRAINT fk_company_id FOREIGN KEY (company_id) REFERENCES companies(id);
ALTER TABLE group_members ADD CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(id);
ALTER TABLE group_members ADD CONSTRAINT fk_idol_id FOREIGN KEY (idol_id) REFERENCES idols(id);
ALTER TABLE idol_info ADD CONSTRAINT fk_idol_id FOREIGN KEY (idol_id) REFERENCES idols(id);

-- Unique constraints
ALTER TABLE group_members ADD CONSTRAINT unique_group_idol UNIQUE (group_id, idol_id);
ALTER TABLE idol_info ADD CONSTRAINT unique_idol_info UNIQUE (idol_id);
ALTER TABLE companies ADD CONSTRAINT unique_company_name UNIQUE (name, country);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE companies CASCADE;

DROP TABLE groups CASCADE;
DROP TABLE group_members;

DROP TABLE idols CASCADE;
DROP TABLE idol_info;

-- +goose StatementEnd
