package handlers

import "strconv"

func parseID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
