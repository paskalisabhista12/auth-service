package utils

import (
	"auth-service/pkg/utils/exception"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	mu          sync.Mutex
	lastSecond  string
	sequenceNum int
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

func GenerateTransactionID() string {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	sec := now.Format("20060102150405") // YYYYMMDDhhmmss

	if sec == lastSecond {
		sequenceNum = (sequenceNum + 1) % 100 // 00–99
	} else {
		lastSecond = sec
		sequenceNum = 0
	}

	return fmt.Sprintf("TRX%s%02d", sec, sequenceNum)
}
