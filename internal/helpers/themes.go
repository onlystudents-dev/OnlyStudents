package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

	maxThemeNameLength = 50
	maxCSSValueLength  = 64
	maxColorsPerTheme  = 64
)

type Theme struct {
	Name   string            `json:"name"`
	Colors map[string]string `json:"colors"`
}

var (
	themesStyleRegex = regexp.MustCompile(`(?si)(<style\b[^>]*\bid=["']themes["'][^>]*>)(.*?)(</style>)`)
	styleTagRegex    = regexp.MustCompile(`(?si)<style[^>]*>(.*?)</style>`)
	rootBlockRegex   = regexp.MustCompile(`(?s):root\s*\{([^}]+)\}`)
	cssVarRegex      = regexp.MustCompile(`--([a-zA-Z0-9_-]+)\s*:`)

	cssCustomPropertyRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	themeNameRegex         = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9 _-]*$`)
	cssHexColorRegex       = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	cssFuncColorRegex      = regexp.MustCompile(`(?i)^(?:rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch|color)\([0-9a-zA-Z%\s,./+-]+\)$`)
	cssVarReferenceRegex   = regexp.MustCompile(`(?i)^var\(--[a-zA-Z][a-zA-Z0-9_-]*\)$`)

	cssColorKeywords = buildCssColorKeywords()
)

const cssColorKeywordList = `aliceblue antiquewhite aqua aquamarine azure beige bisque black ` +
	`blanchedalmond blue blueviolet brown burlywood cadetblue chartreuse chocolate coral ` +
	`cornflowerblue cornsilk crimson cyan darkblue darkcyan darkgoldenrod darkgray darkgreen ` +
	`darkgrey darkkhaki darkmagenta darkolivegreen darkorange darkorchid darkred darksalmon ` +
	`darkseagreen darkslateblue darkslategray darkslategrey darkturquoise darkviolet deeppink ` +
	`deepskyblue dimgray dimgrey dodgerblue firebrick floralwhite forestgreen fuchsia gainsboro ` +
	`ghostwhite gold goldenrod gray green greenyellow grey honeydew hotpink indianred indigo ` +
	`ivory khaki lavender lavenderblush lawngreen lemonchiffon lightblue lightcoral lightcyan ` +
	`lightgoldenrodyellow lightgray lightgreen lightgrey lightpink lightsalmon lightseagreen ` +
	`lightskyblue lightslategray lightslategrey lightsteelblue lightyellow lime limegreen linen ` +
	`magenta maroon mediumaquamarine mediumblue mediumorchid mediumpurple mediumseagreen ` +
	`mediumslateblue mediumspringgreen mediumturquoise mediumvioletred midnightblue mintcream ` +
	`mistyrose moccasin navajowhite navy oldlace olive olivedrab orange orangered orchid ` +
	`palegoldenrod palegreen paleturquoise palevioletred papayawhip peachpuff peru pink plum ` +
	`powderblue purple rebeccapurple red rosybrown royalblue saddlebrown salmon sandybrown ` +
	`seagreen seashell sienna silver skyblue slateblue slategray slategrey snow springgreen ` +
	`steelblue tan teal thistle tomato turquoise violet wheat white whitesmoke yellow yellowgreen ` +
	`transparent currentcolor inherit initial unset revert`

func buildCssColorKeywords() map[string]struct{} {
	keywords := make(map[string]struct{})
	for _, keyword := range strings.Fields(cssColorKeywordList) {
		keywords[keyword] = struct{}{}
	}
	return keywords
}

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

func normalizeColors(colors map[string]string) map[string]string {
	normalized := make(map[string]string, len(colors))
	for key, val := range colors {
		cleanKey := strings.TrimSpace(strings.TrimPrefix(key, "--"))
		normalized[cleanKey] = strings.TrimSpace(val)
	}
	return normalized
}

func isValidCSSColor(value string) bool {
	if len(value) == 0 || len(value) > maxCSSValueLength {
		return false
	}
	if cssHexColorRegex.MatchString(value) {
		return true
	}
	if cssFuncColorRegex.MatchString(value) {
		return true
	}
	if cssVarReferenceRegex.MatchString(value) {
		return true
	}
	_, ok := cssColorKeywords[strings.ToLower(value)]
	return ok
}

func validateThemeStatic(t Theme) error {
	name := strings.TrimSpace(t.Name)
	if name == "" {
		return fmt.Errorf("theme name cannot be empty")
	}
	if len(name) > maxThemeNameLength {
		return fmt.Errorf("theme name exceeds %d characters", maxThemeNameLength)
	}
	if !themeNameRegex.MatchString(name) {
		return fmt.Errorf("theme name contains invalid characters: %q", name)
	}

	if len(t.Colors) == 0 {
		return fmt.Errorf("theme %q has no colors", name)
	}
	if len(t.Colors) > maxColorsPerTheme {
		return fmt.Errorf("theme %q has more than %d color keys", name, maxColorsPerTheme)
	}

	normalized := normalizeColors(t.Colors)

	var invalidKeys, invalidValues []string
	for key, val := range normalized {
		if !cssCustomPropertyRegex.MatchString(key) {
			invalidKeys = append(invalidKeys, key)
			continue
		}
		if !isValidCSSColor(val) {
			invalidValues = append(invalidValues, fmt.Sprintf("%s=%q", key, val))
		}
	}

	if len(invalidKeys) > 0 {
		return fmt.Errorf("invalid CSS custom property names: %s", strings.Join(invalidKeys, ", "))
	}
	if len(invalidValues) > 0 {
		return fmt.Errorf("invalid CSS color values: %s", strings.Join(invalidValues, ", "))
	}

	return nil
}

func (t Theme) IsValid(ctx context.Context, rdb *redis.Client) error {
	if err := validateThemeStatic(t); err != nil {
		return err
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

	normalized := normalizeColors(t.Colors)

	var missing []string
	for _, key := range requiredKeys {
		if _, exists := normalized[key]; !exists {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required color keys: %s", strings.Join(missing, ", "))
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
		slog.Error("set theme", "err", err.Error())
		return err
	}

	queries := db_queries.New(pool)
	colorsBytes, err := json.Marshal(colors)
	if err != nil {
		slog.Error("set theme", "err", err.Error())
		return err
	}

	_, err = queries.UpsertTheme(ctx, db_queries.UpsertThemeParams{
		Name:   name,
		Colors: colorsBytes,
	})
	if err != nil {
		slog.Error("set theme", "err", err.Error())
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
		if err := validateThemeStatic(theme); err != nil {
			return fmt.Errorf("refusing to build css for invalid theme: %w", err)
		}
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
		slog.Error("refresh css", "err", err.Error())
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

	matches := themesStyleRegex.FindSubmatchIndex(html)
	if matches == nil {
		return html
	}

	result := make([]byte, 0, len(html)+len(css))

	result = append(result, html[:matches[3]]...)

	result = append(result, css...)

	result = append(result, html[matches[4]:]...)

	return result
}
