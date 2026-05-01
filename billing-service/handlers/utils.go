package handlers

import "strconv"

func parseID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func coalesceString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
