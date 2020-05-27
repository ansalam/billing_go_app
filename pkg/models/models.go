package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

// Counters struct holds the request count & page count
type Counters struct {
	AuthenticatorID string `json:"AuthenticatorID"`
	RequestCount    int    `json:"RequestCount"`
	PageCount       int    `json:"PageCount"`
}
