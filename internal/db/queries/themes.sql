-- name: GetThemeByName :one
SELECT name, colors
FROM themes
WHERE name = $1;

-- name: ListThemes :many
SELECT name, colors
FROM themes
ORDER BY name ASC;

-- name: UpsertTheme :one
INSERT INTO themes (name, colors)
VALUES ($1, $2)
    ON CONFLICT (name)
DO UPDATE SET colors = EXCLUDED.colors
           RETURNING name, colors;

-- name: DeleteTheme :execrows
DELETE FROM themes
WHERE name = $1;