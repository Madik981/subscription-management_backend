package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) ensurePlanExists(planID uint) error {
	if h.billingBaseURL == "" {
		return errors.New("billing service url is not configured")
	}

	resp, err := h.httpClient.R().
		SetPathParam("id", strconv.FormatUint(uint64(planID), 10)).
		Get(h.billingBaseURL + "/internal/plans/{id}")
	if err != nil {
		return fmt.Errorf("billing service request failed: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return errors.New("plan not found")
	}
	if resp.IsError() {
		return fmt.Errorf("billing service error: %s", resp.Status())
	}

	return nil
}
