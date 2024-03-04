package parse

import (
	"fmt"
	"strconv"
)

type SetFunc func(string, string) error

func Setup(m map[string]any, set SetFunc) error {
	return setup("", m, set)
}

func setup(key string, value any, set SetFunc) error {
	if value == nil {
		return nil
	}

	switch x := value.(type) {
	case map[string]any:
		for k, v := range x {
			err := setup(join(key, k), v, set)
			if err != nil {
				return err
			}
		}
	case []any:
		for _, v := range x {
			err := setup(key, v, set)
			if err != nil {
				return err
			}
		}
	default:
		val, err := format(x)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %w", key, err)
		}

		err = set(key, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func join(a, b string) string {
	if a == "" {
		return b
	}

	return a + "." + b
}

func format(x any) (string, error) {
	switch v := x.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("couldn't convert %T to string", x)
	}
}
