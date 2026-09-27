package helpers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

const RoutesGroupMWKey = "__mw"

var httpMethods = map[string]struct{}{
	fiber.MethodGet: {}, fiber.MethodPost: {}, fiber.MethodPut: {},
	fiber.MethodPatch: {}, fiber.MethodDelete: {}, fiber.MethodHead: {},
	fiber.MethodOptions: {},
}

func Get(v any) fiber.Map    { return fiber.Map{fiber.MethodGet: v} }
func Post(v any) fiber.Map   { return fiber.Map{fiber.MethodPost: v} }
func Patch(v any) fiber.Map  { return fiber.Map{fiber.MethodPatch: v} }
func Delete(v any) fiber.Map { return fiber.Map{fiber.MethodDelete: v} }

func toHandlers(v any) []fiber.Handler {
	switch t := v.(type) {
	case fiber.Handler:
		return []fiber.Handler{t}
	case []fiber.Handler:
		return t
	default:
		panic(fmt.Sprintf("routes: unsupported handler value %T", v))
	}
}

func concat(a, b []fiber.Handler) []fiber.Handler {
	out := make([]fiber.Handler, 0, len(a)+len(b))
	out = append(out, a...)
	return append(out, b...)
}

func joinPath(prefix, seg string) string {
	if prefix == "" {
		return "/" + seg
	}
	return prefix + "/" + seg
}

func pathOrRoot(prefix string) string {
	if prefix == "" {
		return "/"
	}
	return prefix
}

func isMethod(key string) bool {
	_, ok := httpMethods[key]
	return ok
}

func anyHandlers(hs []fiber.Handler) []any {
	out := make([]any, len(hs))
	for i, h := range hs {
		out[i] = h
	}
	return out
}

func RegisterRoutes(r fiber.Router, node fiber.Map) {
	register(r, node, "", nil)
}

func register(r fiber.Router, node fiber.Map, prefix string, inherited []fiber.Handler) {
	mw := inherited
	if v, ok := node[RoutesGroupMWKey]; ok {
		mw = concat(inherited, toHandlers(v))
	}
	for key, val := range node {
		switch {
		case key == RoutesGroupMWKey:
			continue
		case isMethod(key):
			chain := anyHandlers(concat(mw, toHandlers(val)))
			r.Add([]string{key}, pathOrRoot(prefix), chain[0], chain[1:]...)
		default:
			child, ok := val.(fiber.Map)
			if !ok {
				panic(fmt.Sprintf("routes: %q must be a fiber.Map, got %T", key, val))
			}
			register(r, child, joinPath(prefix, key), mw)
		}
	}
}
