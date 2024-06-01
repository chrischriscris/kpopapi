-- +goose Up
-- +goose StatementBegin

CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    metadata_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE image_metadata (
    id SERIAL PRIMARY KEY,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    landscape BOOLEAN NOT NULL, -- True if width >= height
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE idol_images (
    id SERIAL PRIMARY KEY,
    idol_id INTEGER NOT NULL,
    image_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_images (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,
    image_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Foreign key constraints
ALTER TABLE images ADD CONSTRAINT fk_metadata_id FOREIGN KEY (metadata_id) REFERENCES image_metadata(id);
ALTER TABLE idol_images ADD CONSTRAINT fk_idol_id FOREIGN KEY (idol_id) REFERENCES idols(id);
ALTER TABLE idol_images ADD CONSTRAINT fk_image_id FOREIGN KEY (image_id) REFERENCES images(id);
ALTER TABLE group_images ADD CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(id);
ALTER TABLE group_images ADD CONSTRAINT fk_image_id FOREIGN KEY (image_id) REFERENCES images(id);

-- Unique constraints
ALTER TABLE idol_images ADD CONSTRAINT unique_idol_image UNIQUE (idol_id, image_id);
ALTER TABLE group_images ADD CONSTRAINT unique_group_image UNIQUE (group_id, image_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE images;
DROP TABLE image_metadata;

DROP TABLE idol_images;
DROP TABLE group_images;

-- +goose StatementEnd
