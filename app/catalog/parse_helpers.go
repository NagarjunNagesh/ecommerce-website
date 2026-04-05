package catalog

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
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

// parsePositiveDecimal parses a positive decimal query param from q.
// If the param is absent or empty it returns nil.
func parsePositiveDecimal(q url.Values, key string) (*decimal.Decimal, error) {
	raw := strings.TrimSpace(q.Get(key))
	if raw == "" {
		return nil, nil
	}

	value, err := decimal.NewFromString(raw)
	if err != nil {
		return nil, fmt.Errorf("%s must be a valid number", key)
	}
	if value.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("%s must be greater than 0", key)
	}

	return &value, nil
}
