-- name: GetImageByUrl :one
SELECT * FROM images
WHERE url = $1
LIMIT 1;

-- name: AddImage :one
INSERT INTO images (url) VALUES ($1)
RETURNING *;

-- name: AddImageMetadata :one
INSERT INTO image_metadata (
    image_id,
    width,
    height
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: AddIdolImage :one
INSERT INTO idol_images (
    image_id,
    idol_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: AddGroupImage :one
INSERT INTO group_images (
    image_id,
    group_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetRandomImage :one
SELECT * FROM images
ORDER BY random()
LIMIT 1;

-- name: GetNumberOfImages :one
SELECT COUNT(*) FROM images;
