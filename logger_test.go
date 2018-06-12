package main

import (
	"errors"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	logger := NewLogger()
	logger.Info("msg", "ctx")
	logger.Error("error", errors.New("test error"))
}
