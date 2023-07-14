package models

import "time"

type Request struct {
	Url    string        `json:"url"`
	Alias  string        `json:"alias"`
	Expiry time.Duration `json:"expiry"`
}

type Response struct {
	Url                 string        `json:"url"`
	CustomShort         string        `json:"custom_short"`
	Expiry              time.Duration `json:"expiry"`
	XrateLimitRemaining int           `json:"x-rate-limit-remaining"`
	XrateLimitReset     time.Duration `json:"x-rate-limit-reset"`
}
