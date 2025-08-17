package utils

import (
	"auth-service/pkg/utils/exception"
	"encoding/json"
)

func UnmarshalDynamic[T any](data []byte, key string) (T, error) {
	var all map[string]json.RawMessage
	var result T

	if err := json.Unmarshal(data, &all); err != nil {
		return result, err
	}

	raw, ok := all[key]
	if !ok {
		return result, exception.ErrNotFound
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return result, err
	}

	return result, nil
}
