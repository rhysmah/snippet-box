package models

import "errors"

// Created so handlers aren't concerned with the underlying
// datastore or reliant on datastore-specific errors.
var ErrNoRecord = errors.New("models: no matching record found")
