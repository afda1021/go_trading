package config

import "time"

var Durations map[string]time.Duration

func init() {
	Durations = map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}
}
