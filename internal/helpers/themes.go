package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	redisKeyDynamicCSS    = "theme:dynamic_css"
	redisKeyRootColorKeys = "theme:root_color_keys"
)

type Theme struct {
	Name   string            `json:"name"`
	Colors map[string]string `json:"colors"`
}

var (
	styleTagRegex  = regexp.MustCompile(`(?si)<style[^>]*>(.*?)</style>`)
	rootBlockRegex = regexp.MustCompile(`(?s):root\s*\{([^}]+)\}`)
	cssVarRegex    = regexp.MustCompile(`--([a-zA-Z0-9_-]+)\s*:`)
)

func parseRootKeys(html []byte) []string {
	styleMatches := styleTagRegex.FindAllSubmatch(html, -1)
	if len(styleMatches) == 0 {
		return nil
	}

	var rootBody string
	for _, match := range styleMatches {
		if rootMatch := rootBlockRegex.FindSubmatch(match[1]); len(rootMatch) > 1 {
			rootBody = string(rootMatch[1])
			break
		}
	}
	if rootBody == "" {
		return nil
	}

	varMatches := cssVarRegex.FindAllStringSubmatch(rootBody, -1)
	keys := make([]string, 0, len(varMatches))
	seen := make(map[string]struct{}, len(varMatches))

	for _, m := range varMatches {
		if len(m) > 1 {
			k := m[1]
			if _, exists := seen[k]; !exists {
				seen[k] = struct{}{}
				keys = append(keys, k)
			}
		}
	}
	return keys
}

func (t Theme) IsValid(ctx context.Context, rdb *redis.Client) error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("theme name cannot be empty")
	}

	data, err := rdb.Get(ctx, redisKeyRootColorKeys).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("root color keys are not initialized")
		}
		return fmt.Errorf("failed to fetch root color keys from redis: %w", err)
	}

	var requiredKeys []string
	if err := json.Unmarshal(data, &requiredKeys); err != nil {
		return fmt.Errorf("failed to decode root color keys: %w", err)
	}

	if len(requiredKeys) == 0 {
		return fmt.Errorf("root color keys are empty")
	}

	normalized := make(map[string]string, len(t.Colors))
	for k, v := range t.Colors {
		cleanKey := strings.TrimSpace(strings.TrimPrefix(k, "--"))
		normalized[cleanKey] = strings.TrimSpace(v)
	}

	var missing, empty []string
	for _, key := range requiredKeys {
		val, exists := normalized[key]
		if !exists {
			missing = append(missing, key)
			continue
		}
		if val == "" {
			empty = append(empty, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required color keys: %s", strings.Join(missing, ", "))
	}
	if len(empty) > 0 {
		return fmt.Errorf("empty values for color keys: %s", strings.Join(empty, ", "))
	}

	return nil
}

func GetThemes(ctx context.Context, pool *pgxpool.Pool) ([]Theme, error) {
	queries := db_queries.New(pool)
	rows, err := queries.ListThemes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query themes: %w", err)
	}

	result := make([]Theme, 0, len(rows))
	for _, row := range rows {
		var colors map[string]string
		if err := json.Unmarshal(row.Colors, &colors); err != nil {
			return nil, fmt.Errorf("failed to unmarshal colors for theme %q: %w", row.Name, err)
		}
		result = append(result, Theme{
			Name:   row.Name,
			Colors: colors,
		})
	}
	return result, nil
}

func SetTheme(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, name string, colors map[string]string) error {
	ctx := c.Context()
	theme := Theme{Name: name, Colors: colors}
	if err := theme.IsValid(ctx, rdb); err != nil {
		return err
	}

	queries := db_queries.New(pool)
	colorsBytes, err := json.Marshal(colors)
	if err != nil {
		return err
	}

	_, err = queries.UpdateTheme(ctx, db_queries.UpdateThemeParams{
		Name:   name,
		Colors: colorsBytes,
	})
	if err != nil {
		return err
	}

	return RefreshCSS(ctx, pool, rdb)
}

func DeleteTheme(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, name string) error {
	ctx := c.Context()
	queries := db_queries.New(pool)
	if _, err := queries.DeleteTheme(ctx, name); err != nil {
		return err
	}
	return RefreshCSS(ctx, pool, rdb)
}

func rebuildCSS(ctx context.Context, rdb *redis.Client, themes []Theme) error {
	var builder strings.Builder
	for _, theme := range themes {
		builder.WriteString(fmt.Sprintf("[data-theme=%q]{", theme.Name))
		for key, val := range theme.Colors {
			if !strings.HasPrefix(key, "--") {
				key = "--" + key
			}
			builder.WriteString(fmt.Sprintf("%s:%s;", key, val))
		}
		builder.WriteString("}")
	}

	return rdb.Set(ctx, redisKeyDynamicCSS, builder.String(), 0).Err()
}

func RefreshCSS(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client) error {
	themes, err := GetThemes(ctx, pool)
	if err != nil {
		return err
	}
	return rebuildCSS(ctx, rdb, themes)
}

func InitCss(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client, indexHTML []byte) error {
	keys := parseRootKeys(indexHTML)
	if len(keys) == 0 {
		return fmt.Errorf("failed to parse :root color variables from index.html <style>")
	}

	keysJSON, err := json.Marshal(keys)
	if err != nil {
		return fmt.Errorf("failed to marshal root color keys: %w", err)
	}

	if err := rdb.Set(ctx, redisKeyRootColorKeys, keysJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to persist root color keys to redis: %w", err)
	}

	return RefreshCSS(ctx, pool, rdb)
}

func InjectThemes(ctx context.Context, rdb *redis.Client, html []byte) []byte {
	css, err := rdb.Get(ctx, redisKeyDynamicCSS).Bytes()
	if err != nil || len(css) == 0 {
		return html
	}

	closingStyle := []byte("</style>")
	replacement := append(append([]byte{}, css...), closingStyle...)
	return bytes.Replace(html, closingStyle, replacement, 1)
}
