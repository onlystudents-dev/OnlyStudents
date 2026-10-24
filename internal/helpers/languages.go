package helpers

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

var languages = make(map[string][]byte)

func CacheLanguages() {
	for _, lang := range Languages {
		if lang == "" {
			continue
		}
		language, err := os.ReadFile("frontend/dist/assets/languages/" + lang + ".json") // #nosec G304 -- It's a hardcoded list + the files are public anyway
		if err != nil {
			panic(err)
		}
		languages[lang] = language
	}
}

func GetLanguage(c fiber.Ctx, stored string, auto bool) ([]byte, string) {
	if !auto && stored != "" {
		if language, ok := languages[stored]; ok {
			return language, stored
		}
	}

	type preference struct {
		lang string
		q    float64
	}

	var preferences []preference

	for part := range strings.SplitSeq(c.Get("Accept-Language"), ",") {
		fields := strings.Split(strings.TrimSpace(part), ";")
		lang := strings.TrimSpace(fields[0])
		q := 1.0

		for _, field := range fields[1:] {
			field = strings.TrimSpace(field)
			if after, ok := strings.CutPrefix(field, "q="); ok {
				if value, err := strconv.ParseFloat(after, 64); err == nil {
					q = value
				}
			}
		}

		if lang != "" && q > 0 {
			preferences = append(preferences, preference{lang, q})
		}
	}

	sort.SliceStable(preferences, func(i, j int) bool {
		return preferences[i].q > preferences[j].q
	})

	for _, pref := range preferences {
		if language, ok := languages[pref.lang]; ok {
			return language, pref.lang
		}
	}

	return languages["en-US"], "en-US"
}
