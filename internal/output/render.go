package output

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func PrintSection(title string) {
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", len(title)))
}

func PrintList(items []string) {
	for _, item := range items {
		if strings.TrimSpace(item) == "" {
			continue
		}
		fmt.Printf("- %s\n", item)
	}
}

func String(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(t)
	case fmt.Stringer:
		return strings.TrimSpace(t.String())
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func StringFromMap(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := m[key]; ok {
			if text := String(value); text != "" {
				return text
			}
		}
	}
	return ""
}

func IntFromMap(m map[string]any, keys ...string) int {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		switch t := value.(type) {
		case int:
			return t
		case int64:
			return int(t)
		case float64:
			return int(t)
		case string:
			if n, err := strconv.Atoi(strings.TrimSpace(t)); err == nil {
				return n
			}
		}
	}
	return 0
}

func BoolFromMap(m map[string]any, keys ...string) bool {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		switch t := value.(type) {
		case bool:
			return t
		case string:
			return strings.EqualFold(strings.TrimSpace(t), "true")
		}
	}
	return false
}


func MapFromMap(m map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		if nested, ok := value.(map[string]any); ok {
			return nested
		}
	}
	return nil
}

func SliceFromMap(m map[string]any, keys ...string) []map[string]any {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		items, ok := value.([]any)
		if !ok {
			continue
		}
		out := make([]map[string]any, 0, len(items))
		for _, item := range items {
			if entry, ok := item.(map[string]any); ok {
				out = append(out, entry)
			}
		}
		return out
	}
	return nil
}

func Keys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func StringFromNestedMap(m map[string]any, paths ...[]string) string {
	for _, path := range paths {
		current := any(m)
		valid := true
		for _, key := range path {
			next, ok := current.(map[string]any)
			if !ok {
				valid = false
				break
			}
			current, ok = next[key]
			if !ok || current == nil {
				valid = false
				break
			}
		}
		if valid {
			if text := String(current); text != "" {
				return text
			}
		}
	}
	return ""
}

func JoinStrings(m map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}
		switch t := value.(type) {
		case []any:
			parts := make([]string, 0, len(t))
			for _, item := range t {
				if entry, ok := item.(map[string]any); ok {
					text := FirstNonEmpty(
						StringFromMap(entry, "name", "title", "slug", "id"),
						StringFromNestedMap(
							entry,
							[]string{"topic", "name"},
							[]string{"topic", "title"},
							[]string{"node", "name"},
							[]string{"node", "title"},
							[]string{"attributes", "name"},
							[]string{"attributes", "title"},
						),
					)
					if text != "" {
						parts = append(parts, text)
					}
					continue
				}
				if text := String(item); text != "" {
					parts = append(parts, text)
				}
			}
			if len(parts) > 0 {
				return strings.Join(parts, ", ")
			}
		case []string:
			if len(t) > 0 {
				return strings.Join(t, ", ")
			}
		}
	}
	return ""
}
