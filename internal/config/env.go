package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readString(key string, fallback *string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		if fallback == nil {
			return "", fmt.Errorf("missing environment variable %q", key)
		} else {
			return *fallback, nil
		}
	}

	return strings.TrimSpace(val), nil
}

func readInt(key string, fallback *int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		if fallback == nil {
			return 0, fmt.Errorf("missing environment variable %q", key)
		} else {
			return *fallback, nil
		}
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("value for %q (i.e. %q) is not a valid int", key, val)
	}

	return parsed, nil
}
