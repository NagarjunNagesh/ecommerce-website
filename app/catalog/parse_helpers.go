package catalog

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// parseIntOrDefault parses integer query param `key` from `q`.
// If the param is absent or empty it returns `def` as in default.
// On parse error it returns an error with message "<key> must be a valid integer".
func parseIntOrDefault(q url.Values, key string, def int) (int, error) {
	raw := strings.TrimSpace(q.Get(key))
	if raw == "" {
		return def, nil
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", key)
	}
	return v, nil
}
