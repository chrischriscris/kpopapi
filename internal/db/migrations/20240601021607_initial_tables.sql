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
    type CHAR(2) NOT NULL, -- GG, BG, CE (Co-ed)
    debut_date DATE,
    company_id INTEGER REFERENCES companies(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    idol_id INTEGER REFERENCES idols(id),
    since_date DATE,
    until_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE idols (
    id SERIAL PRIMARY KEY,
    stage_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    gender CHAR NOT NULL,
    idol_info_id INTEGER REFERENCES idol_info(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE idol_info (
    id SERIAL PRIMARY KEY,
    birthdate DATE,
    height_cm FLOAT,
    weight_kg FLOAT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE companies;

DROP TABLE groups;
DROP TABLE group_members;

DROP TABLE idols;
DROP TABLE idol_info;

-- +goose StatementEnd
