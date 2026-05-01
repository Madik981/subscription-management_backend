package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) ensureUserExists(userID uint) error {
	if h.accountsBaseURL == "" {
		return errors.New("accounts service url is not configured")
	}

	resp, err := h.httpClient.R().
		SetPathParam("id", strconv.FormatUint(uint64(userID), 10)).
		Get(h.accountsBaseURL + "/internal/users/{id}")
	if err != nil {
		return fmt.Errorf("accounts service request failed: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return errors.New("user not found")
	}
	if resp.IsError() {
		return fmt.Errorf("accounts service error: %s", resp.Status())
	}

	return nil
}
