package helper

import (
	"errors"
	"strings"
)

func InPath(data interface{}, path string) (interface{}, error) {
	for _, p := range strings.Split(path, ".") {
		ok := true

		if p != "" {
			data, ok = data.(map[string]interface{})[p]
		}

		if !ok {
			return nil, errors.New("No data found by path " + path)
		}
	}

	return data, nil
}
