package main

import "testing"

func TestSetUpLog(t *testing.T) {
	cfg := Config{
		LogLevel: "INFO",
		LogType:  "json",
	}

	err := setUpLog(cfg)
	if err != nil {
		t.Errorf("setUpLog() failed: %v", err)
	}

	cfg.LogLevel = "INVALID"

	err = setUpLog(cfg)
	if err == nil {
		t.Errorf("setUpLog() should have failed")
	}

	cfg.LogLevel = "INVALID"
	err = setUpLog(cfg)
	if err == nil {
		t.Errorf("setUpLog() should have failed")
	}
}
