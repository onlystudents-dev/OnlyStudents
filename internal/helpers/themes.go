package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Theme struct {
	Name   string            `json:"name"`
	Colors map[string]string `json:"colors"`
}

var (
	mu            sync.RWMutex
	dynamicCSS    []byte
	rootColorKeys []string

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

func (t Theme) IsValid() error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("theme name cannot be empty")
	}

	mu.RLock()
	requiredKeys := rootColorKeys
	mu.RUnlock()

	if len(requiredKeys) == 0 {
		return fmt.Errorf("root color keys are not initialized")
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

func SetTheme(c fiber.Ctx, pool *pgxpool.Pool, name string, colors map[string]string) error {
	theme := Theme{Name: name, Colors: colors}
	if err := theme.IsValid(); err != nil {
		return err
	}

	queries := db_queries.New(pool)
	colorsBytes, err := json.Marshal(colors)
	if err != nil {
		return err
	}

	_, err = queries.UpdateTheme(c.Context(), db_queries.UpdateThemeParams{
		Name:   name,
		Colors: colorsBytes,
	})
	if err != nil {
		return err
	}

	return RefreshCSS(c.Context(), pool)
}

func DeleteTheme(c fiber.Ctx, pool *pgxpool.Pool, name string) error {
	queries := db_queries.New(pool)
	if _, err := queries.DeleteTheme(c.Context(), name); err != nil {
		return err
	}
	return RefreshCSS(c.Context(), pool)
}

func rebuildCSS(themes []Theme) {
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

	mu.Lock()
	dynamicCSS = []byte(builder.String())
	mu.Unlock()
}

func RefreshCSS(ctx context.Context, pool *pgxpool.Pool) error {
	themes, err := GetThemes(ctx, pool)
	if err != nil {
		return err
	}
	rebuildCSS(themes)
	return nil
}

func InitCss(ctx context.Context, pool *pgxpool.Pool, indexHTML []byte) error {
	keys := parseRootKeys(indexHTML)
	if len(keys) == 0 {
		return fmt.Errorf("failed to parse :root color variables from index.html <style>")
	}

	mu.Lock()
	rootColorKeys = keys
	mu.Unlock()

	return RefreshCSS(ctx, pool)
}

func InjectThemes(html []byte) []byte {
	mu.RLock()
	css := dynamicCSS
	mu.RUnlock()

	if len(css) == 0 {
		return html
	}

	closingStyle := []byte("</style>")
	replacement := append(append([]byte{}, css...), closingStyle...)
	return bytes.Replace(html, closingStyle, replacement, 1)
}
